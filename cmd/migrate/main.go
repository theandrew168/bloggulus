package main

import (
	"context"
	"embed"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v4"
)

//go:embed migrations
var migrationsFS embed.FS

func main() {
	databaseURL := os.Getenv("BLOGGULUS_DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Missing required env var: BLOGGULUS_DATABASE_URL")
	}

	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	if err = conn.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	if err = migrate(conn, context.Background(), migrationsFS); err != nil {
		log.Fatal(err)
	}
}

func migrate(conn *pgx.Conn, ctx context.Context, migrationFS embed.FS) error {
	// create migrations table if it doesn't exist
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS migration (
			migration_id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)`)
	if err != nil {
		return err
	}

	// get migrations that are already applied
	rows, err := conn.Query(ctx, "SELECT name FROM migration")
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

	// get migrations that should be applied (from migrations/ dir)
	migrations, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	// determine missing migrations
	var missing []string
	for _, migration := range migrations {
		name := migration.Name()
		log.Printf("insert into migration(name) values ('%s');\n", name)
		if _, ok := applied[name]; !ok {
			missing = append(missing, name)
		}
	}

	// sort missing migrations to preserve order
	sort.Strings(missing)
	for _, name := range missing {
		log.Printf("applying: %s\n", name)

		// apply the missing ones
		path := filepath.Join("migrations", name)
		sql, err := migrationsFS.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = conn.Exec(ctx, string(sql))
		if err != nil {
			return err
		}

		// update migrations table
		_, err = conn.Exec(ctx, "INSERT INTO migration (name) VALUES ($1)", name)
		if err != nil {
			return err
		}
	}

	log.Printf("migrations up to date!\n")
	return nil
}
