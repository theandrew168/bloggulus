package test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/postgres"
	adminStorage "github.com/theandrew168/bloggulus/backend/postgres/admin/storage"
)

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
