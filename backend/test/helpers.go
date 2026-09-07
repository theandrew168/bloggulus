package test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
)

type CloserFunc func()

func NewConfig(t *testing.T) config.Config {
	t.Helper()

	// This config should mirror what is in bloggulus.test.conf. Since each
	// test's CWD is different based on where the test file is located, there
	// isn't a consistent way to read the test config file directly.
	cfg := config.Config{
		DatabaseURI: "postgresql://postgres:postgres@localhost:5433/postgres",
		SecretKey:   "Be sure to generate your own long and random secret key!",
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

func NewRepository(t *testing.T) (*repository.Repository, CloserFunc) {
	t.Helper()

	db, closer := NewDatabase(t)
	repo := repository.New(db)
	return repo, closer
}
