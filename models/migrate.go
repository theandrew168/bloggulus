package models

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Migrate(db *pgxpool.Pool, migrationsGlob string) error {
	// create migrations table if it doesn't exist
	_, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS migrations (
			migration_id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)`)
	if err != nil {
		return err
	}

	// get migrations that are already applied
	rows, err := db.Query(context.Background(), "SELECT name FROM migrations")
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
	migrations, err := filepath.Glob(migrationsGlob)
	if err != nil {
		return err
	}

	// determine missing migrations
	var missing []string
	for _, migration := range migrations {
		if _, ok := applied[migration]; !ok {
			missing = append(missing, migration)
		}
	}

	// sort missing migrations to preserve order
	sort.Strings(missing)
	for _, file := range missing {
		fmt.Printf("applying: %s\n", file)

		// apply the missing ones
		sql, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		_, err = db.Exec(context.Background(), string(sql))
		if err != nil {
			return err
		}

		// update migrations table
		_, err = db.Exec(context.Background(), "INSERT INTO migrations (name) VALUES ($1)", file)
		if err != nil {
			return err
		}
	}

	return nil
}
