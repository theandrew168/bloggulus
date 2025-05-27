package test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/feed"
	feedMock "github.com/theandrew168/bloggulus/backend/feed/mock"
	"github.com/theandrew168/bloggulus/backend/job"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/query"
	"github.com/theandrew168/bloggulus/backend/repository"
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

func NewQuery(t *testing.T) (*query.Query, CloserFunc) {
	t.Helper()

	db, closer := NewDatabase(t)
	q := query.New(db)
	return q, closer
}

func NewSyncService(
	t *testing.T,
	repo *repository.Repository,
	feeds map[string]feed.FetchFeedResponse,
) *job.SyncService {
	t.Helper()

	syncService := job.NewSyncService(repo, feedMock.NewFeedFetcher(feeds))
	return syncService
}
