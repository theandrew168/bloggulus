package command

import (
	"context"
	"errors"
	"log/slog"

	"golang.org/x/sync/semaphore"

	"github.com/theandrew168/bloggulus/backend/command/sync"
	"github.com/theandrew168/bloggulus/backend/feed"
	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/timeutil"
	"github.com/theandrew168/bloggulus/backend/value"
)

const (
	// How many blogs to sync at once.
	SyncConcurrency = 4
)

type SyncCommand struct {
	repo        *repository.Repository
	feedFetcher feed.FeedFetcher
}

func NewSync(repo *repository.Repository, feedFetcher feed.FeedFetcher) *SyncCommand {
	cmd := SyncCommand{
		repo:        repo,
		feedFetcher: feedFetcher,
	}
	return &cmd
}

// Sync a new or existing Blog based on the provided feed URL.
func (cmd *SyncCommand) SyncBlog(ctx context.Context, feedURL value.URL) error {
	blog, err := cmd.repo.Blog().ReadByFeedURL(ctx, feedURL)
	if err != nil {
		if !errors.Is(err, postgres.ErrNotFound) {
			return err
		}

		// An ErrNotFound is acceptable (and expected) here. The only difference
		// is that we won't be able to include the ETag and Last-Modified headers
		// in the request. This is fine for new blogs (an unconditional fetch).
		return sync.SyncNewBlog(ctx, cmd.repo, cmd.feedFetcher, feedURL)
	}

	return sync.SyncExistingBlog(ctx, cmd.repo, cmd.feedFetcher, blog)
}

// Start with the current time and a list of all known blogs. For each blog,
// compare its syncedAt time to the current time. If the difference is less
// than SyncCooldown, skip it. Otherwise, check for and sync new content.
func (cmd *SyncCommand) SyncAllBlogs(ctx context.Context) error {
	slog.Info("syncing blogs")

	blogs, err := cmd.repo.Blog().List(ctx)
	if err != nil {
		return err
	}

	// Be sure to only sync blogs that are ready.
	now := timeutil.Now()
	syncableBlogs := sync.FilterSyncableBlogs(blogs, now)

	// Update the syncedAt time for each syncable blog before syncing.
	for _, blog := range syncableBlogs {
		blog.SetSyncedAt(now)
		err = cmd.repo.Blog().Update(ctx, blog)
		if err != nil {
			return err
		}
	}

	ParallelForEach(SyncConcurrency, syncableBlogs, func(blog *model.Blog) {
		slog.Info("syncing blog", "title", blog.Title().Value(), "id", blog.ID().String())
		err = cmd.SyncBlog(ctx, blog.FeedURL())
		if err != nil {
			slog.Warn(err.Error(), "title", blog.Title().Value(), "id", blog.ID().String())
		}
	})

	return nil
}

func ParallelForEach[T any](concurrency int, items []T, fn func(T)) {
	// Use a weighted semaphore to limit concurrency.
	sem := semaphore.NewWeighted(int64(concurrency))

	// Perform tasks in parallel (up to "concurrency" at once).
	for _, item := range items {
		sem.Acquire(context.Background(), 1)

		go func(item T) {
			defer sem.Release(1)
			fn(item)
		}(item)
	}

	// Wait for all tasks to finish.
	sem.Acquire(context.Background(), SyncConcurrency)
}
