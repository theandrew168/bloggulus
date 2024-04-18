package service

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/feed"
	"github.com/theandrew168/bloggulus/backend/fetch"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
)

const SyncInterval = 1 * time.Hour

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
	// perform an initial sync upfront
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
// compare its syncedAt time to the current time. If the difference is an hour
// or larger, sync the blog. Otherwise, skip syncing it.
func (s *SyncService) SyncAllBlogs() error {
	// ensure only one sync happens at a time
	if !s.mu.TryLock() {
		slog.Info("sync already in progress")
		return nil
	}
	defer s.mu.Unlock()

	slog.Info("syncing blogs")

	blogs, err := s.store.Admin().Blog().ListAll()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, blog := range blogs {
		// only sync each blog once an hour (at most) regardless of SyncInterval
		delta := now.Sub(blog.SyncedAt())
		if delta < 1*time.Hour {
			slog.Info("skipping blog", "title", blog.Title(), "id", blog.ID())
			continue
		}

		slog.Info("syncing blog", "title", blog.Title(), "id", blog.ID())
		err = s.SyncBlog(blog.FeedURL())
		if err != nil {
			slog.Warn(err.Error(), "title", blog.Title(), "id", blog.ID())
		}
	}

	return nil
}

func (s *SyncService) SyncBlog(feedURL string) error {
	blog, err := s.store.Admin().Blog().ReadByFeedURL(feedURL)
	if err != nil {
		if !errors.Is(err, postgres.ErrNotFound) {
			return err
		}

		return s.syncNewBlog(feedURL)
	}

	return s.syncExistingBlog(blog)
}

func (s *SyncService) syncNewBlog(feedURL string) error {
	resp, err := s.feedFetcher.FetchFeed(feedURL, "", "")
	if err != nil {
		return err
	}

	// no feed from a new blog (no cache headers) is an error
	if resp.Feed == "" {
		return fetch.ErrUnreachableFeed
	}

	feedBlog, err := feed.Parse(feedURL, resp.Feed)
	if err != nil {
		return err
	}

	feedBlog, err = feed.Hydrate(feedBlog, s.pageFetcher)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	blog := admin.NewBlog(
		feedBlog.FeedURL,
		feedBlog.SiteURL,
		feedBlog.Title,
		resp.ETag,
		resp.LastModified,
		now,
	)
	err = s.store.Admin().Blog().Create(blog)
	if err != nil {
		return err
	}

	for _, feedPost := range feedBlog.Posts {
		err = s.syncPost(blog, feedPost)
		if err != nil {
			slog.Warn(err.Error(), "url", feedPost.URL)
		}
	}

	return nil
}

func (s *SyncService) syncExistingBlog(blog *admin.Blog) error {
	now := time.Now().UTC()
	blog.SetSyncedAt(now)

	err := s.store.Admin().Blog().Update(blog)
	if err != nil {
		return err
	}

	resp, err := s.feedFetcher.FetchFeed(blog.FeedURL(), blog.ETag(), blog.LastModified())
	if err != nil {
		return err
	}

	if resp.Feed == "" {
		slog.Info("no new content", "title", blog.Title(), "id", blog.ID())
		return nil
	}

	if resp.ETag != "" {
		blog.SetETag(resp.ETag)
	}

	if resp.LastModified != "" {
		blog.SetLastModified(resp.LastModified)
	}

	err = s.store.Admin().Blog().Update(blog)
	if err != nil {
		return err
	}

	feedBlog, err := feed.Parse(blog.FeedURL(), resp.Feed)
	if err != nil {
		return err
	}

	feedBlog, err = feed.Hydrate(feedBlog, s.pageFetcher)
	if err != nil {
		return err
	}

	for _, feedPost := range feedBlog.Posts {
		err = s.syncPost(blog, feedPost)
		if err != nil {
			slog.Warn(err.Error(), "url", feedPost.URL)
		}
	}

	return nil
}

func (s *SyncService) syncPost(blog *admin.Blog, feedPost feed.Post) error {
	post, err := s.store.Admin().Post().ReadByURL(feedPost.URL)
	if err != nil {
		if !errors.Is(err, postgres.ErrNotFound) {
			return err
		}

		post := admin.NewPost(
			blog,
			feedPost.URL,
			feedPost.Title,
			feedPost.Content,
			feedPost.PublishedAt,
		)
		return s.store.Admin().Post().Create(post)
	}

	// update the post's content (if available)
	if feedPost.Content != "" {
		post.SetContent(feedPost.Content)
		return s.store.Admin().Post().Update(post)
	}

	return nil
}
