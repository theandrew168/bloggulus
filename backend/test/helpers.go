package test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/fetch"
	fetchMock "github.com/theandrew168/bloggulus/backend/fetch/mock"
	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/service"
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

func NewRepository(t *testing.T) (*repository.Repository, CloserFunc) {
	t.Helper()

	db, closer := NewDatabase(t)
	repo := repository.New(db)
	return repo, closer
}

func NewFinder(t *testing.T) (*finder.Finder, CloserFunc) {
	t.Helper()

	db, closer := NewDatabase(t)
	find := finder.New(db)
	return find, closer
}

func NewSyncService(
	t *testing.T,
	repo *repository.Repository,
	feeds map[string]fetch.FetchFeedResponse,
) *service.SyncService {
	t.Helper()

	syncService := service.NewSyncService(repo, fetchMock.NewFeedFetcher(feeds))
	return syncService
}
