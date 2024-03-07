package test

import (
	"io"
	"log"
	"testing"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/database"
	"github.com/theandrew168/bloggulus/backend/storage"
)

type CloserFunc func()

func NewLogger(t *testing.T) *log.Logger {
	return log.New(io.Discard, "", 0)
}

func NewConfig(t *testing.T) config.Config {
	t.Helper()

	// TODO: is there a better way to handle this? trying to read the
	// conf in the root dir depends on where the reading _file_ is. wacky.
	cfg := config.Config{
		DatabaseURI: "postgresql://postgres:postgres@localhost:5432/postgres",
	}
	return cfg
}

func NewDatabase(t *testing.T) (database.Conn, CloserFunc) {
	t.Helper()

	cfg := NewConfig(t)
	pool, err := database.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		t.Fatal(err)
	}

	return pool, pool.Close
}

func NewStorage(t *testing.T) (*storage.Storage, CloserFunc) {
	t.Helper()

	db, closer := NewDatabase(t)
	store := storage.New(db)
	return store, closer
}
