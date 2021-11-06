package main_test

import (
	"context"
	"embed"
	"io/fs"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/internal/postgresql"
)

//go:embed migrations
var migrationsFS embed.FS

func TestMigrate(t *testing.T) {
	// check for database connection url var
	databaseURL := os.Getenv("BLOGGULUS_DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("Missing required env var: BLOGGULUS_DATABASE_URL")
	}

	// open a database connection pool
	conn, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// test connection to ensure all is well
	if err = conn.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}

	// apply database migrations
	migrations, _ := fs.Sub(migrationsFS, "migrations")
	err = postgresql.Migrate(conn, context.Background(), migrations)
	if err != nil {
		t.Fatal(err)
	}
}
