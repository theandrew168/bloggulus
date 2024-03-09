package migrate

import (
	"context"
	"io/fs"
	"sort"

	"github.com/theandrew168/bloggulus/backend/database"
)

func Migrate(conn database.Conn, files fs.FS) ([]string, error) {
	ctx := context.Background()

	// create migration table if it doesn't exist
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS migration (
			id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)`)
	if err != nil {
		return nil, err
	}

	// get migrations that are already applied
	rows, err := conn.Query(ctx, "SELECT name FROM migration")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existing := make(map[string]bool)
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		existing[name] = true
	}

	// get migrations that should be applied
	subdir, _ := fs.Sub(files, "migrations")
	migrations, err := fs.ReadDir(subdir, ".")
	if err != nil {
		return nil, err
	}

	// determine missing migrations
	var missing []string
	for _, migration := range migrations {
		name := migration.Name()
		if _, ok := existing[name]; !ok {
			missing = append(missing, name)
		}
	}

	// sort missing migrations to preserve order
	sort.Strings(missing)

	// apply each missing migration
	var applied []string
	for _, name := range missing {
		sql, err := fs.ReadFile(subdir, name)
		if err != nil {
			return nil, err
		}

		// apply each migration in a transaction
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return nil, err
		}
		defer tx.Rollback(context.Background())

		_, err = tx.Exec(ctx, string(sql))
		if err != nil {
			return nil, err
		}

		// update migration table
		_, err = tx.Exec(ctx, "INSERT INTO migration (name) VALUES ($1)", name)
		if err != nil {
			return nil, err
		}

		err = tx.Commit(context.Background())
		if err != nil {
			return nil, err
		}

		applied = append(applied, name)
	}

	return applied, nil
}
