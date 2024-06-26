package test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/config"
	fetch "github.com/theandrew168/bloggulus/backend/fetch/mock"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/storage"
)

type CloserFunc func()

func NewConfig(t *testing.T) config.Config {
	t.Helper()

	// TODO: find a consistent way to read "../../bloggulus.test.conf"
	cfg := config.Config{
		DatabaseURI: "postgresql://postgres:postgres@localhost:5433/postgres",
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

func NewStorage(t *testing.T) (*storage.Storage, CloserFunc) {
	t.Helper()

	db, closer := NewDatabase(t)
	store := storage.New(db)
	return store, closer
}

func NewSyncService(
	t *testing.T,
	store *storage.Storage,
	feeds map[string]string,
	pages map[string]string,
) *service.SyncService {
	t.Helper()

	syncService := service.NewSyncService(store, fetch.NewFeedFetcher(feeds), fetch.NewPageFetcher(pages))
	return syncService
}
