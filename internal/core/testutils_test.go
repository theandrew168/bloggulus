package core_test

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomString(n int) string {
	validRunes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	b := make([]byte, n)
	for i := range b {
		b[i] = validRunes[rand.Intn(len(validRunes))]
	}

	return string(b)
}

func connectDB(t *testing.T) *pgxpool.Pool {
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
		t.Fatal(err)
	}

	return conn
}
