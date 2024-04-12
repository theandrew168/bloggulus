package test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/config"
	adminStorage "github.com/theandrew168/bloggulus/backend/domain/admin/storage/postgres"
	"github.com/theandrew168/bloggulus/backend/postgres"
	readerStorage "github.com/theandrew168/bloggulus/backend/postgres/reader/storage"
)

// TODO: do we even need these?

type CloserFunc func()

func NewConfig(t *testing.T) config.Config {
	t.Helper()

	cfg := config.Config{
		DatabaseURI: "postgresql://postgres:postgres@localhost:5432/postgres",
	}
	return cfg
}

func NewDatabase(t *testing.T) (postgres.Conn, CloserFunc) {
	t.Helper()

	cfg := NewConfig(t)
	pool, err := postgres.ConnectPool(cfg.DatabaseURI)
	AssertNilError(t, err)

	return pool, pool.Close
}

func NewAdminStorage(t *testing.T) (*adminStorage.Storage, CloserFunc) {
	t.Helper()

	db, closer := NewDatabase(t)
	store := adminStorage.New(db)
	return store, closer
}

func NewReaderStorage(t *testing.T) (*readerStorage.Storage, CloserFunc) {
	t.Helper()

	db, closer := NewDatabase(t)
	store := readerStorage.New(db)
	return store, closer
}
