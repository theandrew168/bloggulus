package service

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/theandrew168/bloggulus/backend/feed"
	"github.com/theandrew168/bloggulus/backend/fetch"
	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

const (
	// Check for new posts every SyncInterval.
	SyncInterval = 2 * time.Hour
	// Since blogs are synced at startup, ensure they don't get synced more often than every SyncCooldown.
	SyncCooldown = 1 * time.Hour

	// How many blogs to sync at once.
	SyncConcurrency = 8
)

type SyncService struct {
	mu          sync.Mutex
	store       *storage.Storage
	feedFetcher fetch.FeedFetcher
	pageFetcher fetch.PageFetcher
}

func NewSyncService(store *storage.Storage, feedFetcher fetch.FeedFetcher, pageFetcher fetch.PageFetcher) *SyncService {
	s := SyncService{
		store:       store,
		feedFetcher: feedFetcher,
		pageFetcher: pageFetcher,
	}
	return &s
}

func (s *SyncService) Run(ctx context.Context) error {
	// perform an initial sync at service startup
	err := s.SyncAllBlogs()
	if err != nil {
		slog.Error(err.Error())
	}

	// then again every "internal" until stopped
	ticker := time.NewTicker(SyncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("stopping sync service")
			slog.Info("stopped sync service")
			return nil
		case <-ticker.C:
			err := s.SyncAllBlogs()
			if err != nil {
				slog.Error(err.Error())
			}
		}
	}

}

// Start with the current time and a list of all known blogs. For each blog,
// compare its syncedAt time to the current time. If the difference is less
// than SyncCooldown, skip it. Otherwise, check for and sync new content.
func (s *SyncService) SyncAllBlogs() error {
	// ensure only one sync happens at a time
	if !s.mu.TryLock() {
		slog.Info("sync already in progress")
		return nil
	}
	defer s.mu.Unlock()

	slog.Info("syncing blogs")

	blogs, err := s.store.Blog().ListAll()
	if err != nil {
		return err
	}

	now := timeutil.Now()

	// use a weighted semaphore to limit concurrency
	sem := semaphore.NewWeighted(SyncConcurrency)

	// sync all blogs in parallel (up to SyncConcurrency at once)
	for _, blog := range blogs {
		sem.Acquire(context.Background(), 1)

		go func(blog *model.Blog) {
			defer sem.Release(1)

			// don't sync a given blog more than once per SyncCooldown
			delta := now.Sub(blog.SyncedAt())
			if delta < SyncCooldown {
				slog.Info("skipping blog", "title", blog.Title(), "id", blog.ID())
				return
			}

			slog.Info("syncing blog", "title", blog.Title(), "id", blog.ID())
			_, err = s.SyncBlog(blog.FeedURL())
			if err != nil {
				slog.Warn(err.Error(), "title", blog.Title(), "id", blog.ID())
				return
			}
		}(blog)
	}

	// wait for all blogs to finish syncing
	sem.Acquire(context.Background(), SyncConcurrency)

	return nil
}

func (s *SyncService) SyncBlog(feedURL string) (*model.Blog, error) {
	blog, err := s.store.Blog().ReadByFeedURL(feedURL)
	if err != nil {
		if !errors.Is(err, postgres.ErrNotFound) {
			return nil, err
		}

		return s.syncNewBlog(feedURL)
	}

	return s.syncExistingBlog(blog)
}

func (s *SyncService) syncNewBlog(feedURL string) (*model.Blog, error) {
	resp, err := s.feedFetcher.FetchFeed(feedURL, "", "")
	if err != nil {
		return nil, err
	}

	// no feed from a new blog (no cache headers) is an error
	if resp.Feed == "" {
		return nil, fetch.ErrUnreachableFeed
	}

	feedBlog, err := feed.Parse(feedURL, resp.Feed)
	if err != nil {
		return nil, err
	}

	feedBlog, err = feed.Hydrate(feedBlog, s.pageFetcher)
	if err != nil {
		return nil, err
	}

	now := timeutil.Now()
	blog, err := model.NewBlog(
		feedBlog.FeedURL,
		feedBlog.SiteURL,
		feedBlog.Title,
		resp.ETag,
		resp.LastModified,
		now,
	)
	if err != nil {
		return nil, err
	}

	err = s.store.Blog().Create(blog)
	if err != nil {
		return nil, err
	}

	for _, feedPost := range feedBlog.Posts {
		_, err = s.syncPost(blog, feedPost)
		if err != nil {
			slog.Warn(err.Error(), "url", feedPost.URL)
		}
	}

	return blog, nil
}

func (s *SyncService) syncExistingBlog(blog *model.Blog) (*model.Blog, error) {
	now := timeutil.Now()
	blog.SetSyncedAt(now)

	err := s.store.Blog().Update(blog)
	if err != nil {
		return nil, err
	}

	resp, err := s.feedFetcher.FetchFeed(blog.FeedURL(), blog.ETag(), blog.LastModified())
	if err != nil {
		return nil, err
	}

	if resp.Feed == "" {
		slog.Info("no new content", "title", blog.Title(), "id", blog.ID())
		return blog, nil
	}

	if resp.ETag != "" {
		blog.SetETag(resp.ETag)
	}

	if resp.LastModified != "" {
		blog.SetLastModified(resp.LastModified)
	}

	err = s.store.Blog().Update(blog)
	if err != nil {
		return nil, err
	}

	feedBlog, err := feed.Parse(blog.FeedURL(), resp.Feed)
	if err != nil {
		return nil, err
	}

	feedBlog, err = feed.Hydrate(feedBlog, s.pageFetcher)
	if err != nil {
		return nil, err
	}

	for _, feedPost := range feedBlog.Posts {
		_, err = s.syncPost(blog, feedPost)
		if err != nil {
			slog.Warn(err.Error(), "url", feedPost.URL)
		}
	}

	return blog, nil
}

func (s *SyncService) syncPost(blog *model.Blog, feedPost feed.Post) (*model.Post, error) {
	post, err := s.store.Post().ReadByURL(feedPost.URL)
	if err != nil {
		if !errors.Is(err, postgres.ErrNotFound) {
			return nil, err
		}

		post, err := model.NewPost(
			blog,
			feedPost.URL,
			feedPost.Title,
			feedPost.Content,
			feedPost.PublishedAt,
		)
		if err != nil {
			return nil, err
		}

		err = s.store.Post().Create(post)
		if err != nil {
			return nil, err
		}

		return post, nil
	}

	// update the post's content (if available)
	if feedPost.Content != "" {
		post.SetContent(feedPost.Content)
		err = s.store.Post().Update(post)
		if err != nil {
			return nil, err
		}

		return post, nil
	}

	return post, nil
}
