package service

import (
	"errors"
	"log/slog"
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
			slog.Info("recently synced", "title", blog.Title)
		}

		s.SyncBlog(blog.FeedURL)
	}
	return nil
}

func (s *SyncService) SyncBlog(feedURL string) error {
	blog, err := s.store.Blog().ReadByFeedURL(feedURL)
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

	if resp.Feed == "" {
		return errors.New("sync: skipping due to empty / up-to-date feed")
	}

	feedBlog, err := feed.Parse(feedURL, resp.Feed)
	if err != nil {
		return err
	}

	feedBlog, err = feed.Hydrate(feedBlog, s.pageFetcher)
	if err != nil {
		return err
	}

	blog := admin.NewBlog(
		feedBlog.FeedURL,
		feedBlog.SiteURL,
		feedBlog.Title,
		resp.ETag,
		resp.LastModified,
		time.Now().Add(-1*time.Hour),
	)
	err = s.store.Blog().Create(blog)
	if err != nil {
		return err
	}

	for _, feedPost := range feedBlog.Posts {
		err = s.syncPost(blog, feedPost)
		if err != nil {
			slog.Error(err.Error())
		}
	}

	return nil
}

func (s *SyncService) syncExistingBlog(blog admin.Blog) error {
	resp, err := s.feedFetcher.FetchFeed(blog.FeedURL, blog.ETag, blog.LastModified)
	if err != nil {
		return err
	}

	if resp.Feed == "" {
		return errors.New("sync: skipping due to empty / up-to-date feed")
	}

	if resp.ETag != "" {
		blog.ETag = resp.ETag
	}

	if resp.LastModified != "" {
		blog.LastModified = resp.LastModified
	}

	err = s.store.Blog().Update(blog)
	if err != nil {
		return err
	}

	feedBlog, err := feed.Parse(blog.FeedURL, resp.Feed)
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
			slog.Error(err.Error())
		}
	}

	return nil
}

func (s *SyncService) syncPost(blog admin.Blog, feedPost feed.Post) error {
	post, err := s.store.Post().ReadByURL(feedPost.URL)
	if err != nil {
		if !errors.Is(err, postgres.ErrNotFound) {
			return err
		}

		post := admin.NewPost(
			blog,
			feedPost.URL,
			feedPost.Title,
			feedPost.Contents,
			feedPost.PublishedAt,
		)
		return s.store.Post().Create(post)
	}

	// update the post's contents (if available)
	if feedPost.Contents != "" {
		post.Contents = feedPost.Contents
		return s.store.Post().Update(post)
	}

	return nil
}
