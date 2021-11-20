package test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/internal/config"
)

func ConnectDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	// use config defaults for tests
	appConfig := config.Defaults()

	// open a database connection pool
	conn, err := pgxpool.Connect(context.Background(), appConfig.DatabaseURL)
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
