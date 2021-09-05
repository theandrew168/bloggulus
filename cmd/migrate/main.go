package main

import (
	"context"
	"embed"
	"log"
	"os"

	"github.com/jackc/pgx/v4"

	"github.com/theandrew168/bloggulus/internal/postgresql"
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

	if err = postgresql.Migrate(conn, context.Background(), migrationsFS); err != nil {
		log.Fatal(err)
	}
}
