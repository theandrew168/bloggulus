package test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

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

	// test connection to ensure all is well
	if err = conn.Ping(context.Background()); err != nil {
		conn.Close()
		t.Fatal(err)
	}

	return conn
}
