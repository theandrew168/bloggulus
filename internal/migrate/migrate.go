package migrate

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"sort"

	"github.com/theandrew168/bloggulus/internal/database"
)

//go:embed migration
var migrationFS embed.FS

func Migrate(db database.Conn, logger *log.Logger) error {
	ctx := context.Background()

	// create migrations table if it doesn't exist
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS migration (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)`)
	if err != nil {
		return err
	}

	// get migrations that are already applied
	rows, err := db.Query(ctx, "SELECT name FROM migration")
	if err != nil {
		return err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return err
		}
		applied[name] = true
	}

	// get migrations that should be applied (from migration/ dir)
	subdir, _ := fs.Sub(migrationFS, "migration")
	migrations, err := fs.ReadDir(subdir, ".")
	if err != nil {
		return err
	}

	// determine missing migrations
	var missing []string
	for _, migration := range migrations {
		name := migration.Name()
		if _, ok := applied[name]; !ok {
			missing = append(missing, name)
		}
	}

	// sort missing migrations to preserve order
	sort.Strings(missing)
	for _, name := range missing {
		logger.Printf("applying: %s\n", name)

		// apply the missing ones
		sql, err := fs.ReadFile(subdir, name)
		if err != nil {
			return err
		}
		_, err = db.Exec(ctx, string(sql))
		if err != nil {
			return err
		}

		// update migrations table
		_, err = db.Exec(ctx, "INSERT INTO migration (name) VALUES ($1)", name)
		if err != nil {
			return err
		}
	}

	logger.Println("migrations up to date")
	return nil
}
