package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/feed"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

type SyncService struct {
	store       storage.Storage
	feedFetcher feed.FeedFetcher
	pageFetcher feed.PageFetcher
}

func NewSyncService(store storage.Storage, feedFetcher feed.FeedFetcher, pageFetcher feed.PageFetcher) *SyncService {
	s := SyncService{
		store:       store,
		feedFetcher: feedFetcher,
		pageFetcher: pageFetcher,
	}
	return &s
}

// Start with the current time and a list of all known blogs. For each blog,
// compare its syncedAt time to the current time. If the difference is an hour
// or larger, sync the blog. Otherwise, skip syncing it.
func (s *SyncService) SyncAllBlogs() error {
	now := time.Now()

	blogs, err := s.store.Blog().List(1000, 0)
	if err != nil {
		return err
	}

	for _, blog := range blogs {
		delta := now.Sub(blog.SyncedAt)
		if delta < time.Hour {
			fmt.Printf("recently synced: %v\n", blog.Title)
		}

		s.SyncBlog(blog.FeedURL)
	}
	return nil
}

func (s *SyncService) SyncBlog(feedURL string) error {
	blog, err := s.store.Blog().ReadByFeedURL(feedURL)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return s.syncNewBlog(feedURL)
		}

		return err
	}

	return s.syncExistingBlog(blog)
}

func (s *SyncService) syncNewBlog(feedURL string) error {
	return nil
}

func (s *SyncService) syncExistingBlog(blog admin.Blog) error {
	return nil
}
