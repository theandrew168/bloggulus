package test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/fetch"
	fetchMock "github.com/theandrew168/bloggulus/backend/fetch/mock"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/query"
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
	s := storage.New(db)
	return s, closer
}

func NewQuery(t *testing.T) (*query.Query, CloserFunc) {
	t.Helper()

	db, closer := NewDatabase(t)
	q := query.New(db)
	return q, closer
}

func NewSyncService(
	t *testing.T,
	store *storage.Storage,
	feeds map[string]fetch.FetchFeedResponse,
	pages map[string]string,
) *service.SyncService {
	t.Helper()

	syncService := service.NewSyncService(store, fetchMock.NewFeedFetcher(feeds), fetchMock.NewPageFetcher(pages))
	return syncService
}
