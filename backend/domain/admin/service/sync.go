package service

import (
	"errors"
	"log/slog"
	"time"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/feed"
	"github.com/theandrew168/bloggulus/backend/domain/admin/fetch"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
)

type SyncService struct {
	store       storage.Storage
	feedFetcher fetch.FeedFetcher
	pageFetcher fetch.PageFetcher
}

func NewSyncService(store storage.Storage, feedFetcher fetch.FeedFetcher, pageFetcher fetch.PageFetcher) *SyncService {
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
	now := time.Now().UTC()

	// TODO: page through and collect all blogs
	blogs, err := s.store.Blog().List(1000, 0)
	if err != nil {
		return err
	}

	for _, blog := range blogs {
		delta := now.Sub(blog.SyncedAt())
		if delta < time.Hour {
			slog.Info("recently synced", "title", blog.Title)
			continue
		}

		s.SyncBlog(blog.FeedURL())
	}
	return nil
}

func (s *SyncService) SyncBlog(feedURL string) error {
	blog, err := s.store.Blog().ReadByFeedURL(feedURL)
	if err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
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

func (s *SyncService) syncExistingBlog(blog *admin.Blog) error {
	now := time.Now().UTC()
	blog.SetSyncedAt(now)

	err := s.store.Blog().Update(blog)
	if err != nil {
		return err
	}

	resp, err := s.feedFetcher.FetchFeed(blog.FeedURL(), blog.ETag(), blog.LastModified())
	if err != nil {
		return err
	}

	if resp.ETag != "" {
		blog.SetETag(resp.ETag)
	}

	if resp.LastModified != "" {
		blog.SetLastModified(resp.LastModified)
	}

	err = s.store.Blog().Update(blog)
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
			slog.Error(err.Error())
		}
	}

	return nil
}

func (s *SyncService) syncPost(blog *admin.Blog, feedPost feed.Post) error {
	post, err := s.store.Post().ReadByURL(feedPost.URL)
	if err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			return err
		}

		post := admin.NewPost(
			blog,
			feedPost.URL,
			feedPost.Title,
			feedPost.Content,
			feedPost.PublishedAt,
		)
		return s.store.Post().Create(post)
	}

	// update the post's content (if available)
	if feedPost.Content != "" {
		post.SetContent(feedPost.Content)
		return s.store.Post().Update(post)
	}

	return nil
}
