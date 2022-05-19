package test

import (
	"io"
	"log"
	"testing"

	"github.com/theandrew168/bloggulus/internal/config"
	"github.com/theandrew168/bloggulus/internal/database"
	"github.com/theandrew168/bloggulus/internal/storage"
)

type CloserFunc func()

func NewLogger(t *testing.T) *log.Logger {
	return log.New(io.Discard, "", 0)
}

func NewConfig(t *testing.T) config.Config {
	t.Helper()

	// read the local development config file
	cfg, err := config.ReadFile("../../bloggulus.conf")
	if err != nil {
		t.Fatal(err)
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
