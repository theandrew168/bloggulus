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
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

// SYNC:
// Calculation: Run the sync process every N hours (time.Ticker)
// Action: List all blogs in the database
// Calculation: Determine which blogs need to be synced (FilterSyncableBlogs)
// ADD:
// Calculation: For each sync-able blog, create it's FetchFeedRequest (CreateSyncRequest)
// Action: Update sync time for each sync-able blog
// Action: Exchange the request to fetch the blog for a response (with limited concurrency)
// Action: Update the blog's cache headers if changed
// Calculation: If the response includes data, parse the RSS / Atom feed for posts
// Action: List all posts for the current blog
// Calculation: Determine if any posts in the feed are new / updated
// Action: Create / update posts in the database

const (
	// Check for new posts every SyncInterval.
	SyncInterval = 2 * time.Hour

	// How many blogs to sync at once.
	SyncConcurrency = 8
)

func FilterSyncableBlogs(blogs []*model.Blog, now time.Time) []*model.Blog {
	var syncableBlogs []*model.Blog
	for _, blog := range blogs {
		if blog.CanBeSynced(now) {
			syncableBlogs = append(syncableBlogs, blog)
		}
	}
	return syncableBlogs
}

type SyncService struct {
	mu          sync.Mutex
	repo        *repository.Repository
	feedFetcher fetch.FeedFetcher
	pageFetcher fetch.PageFetcher
}

func NewSyncService(repo *repository.Repository, feedFetcher fetch.FeedFetcher, pageFetcher fetch.PageFetcher) *SyncService {
	s := SyncService{
		repo:        repo,
		feedFetcher: feedFetcher,
		pageFetcher: pageFetcher,
	}
	return &s
}

func (s *SyncService) updateCacheHeaders(blog *model.Blog, response fetch.FetchFeedResponse) error {
	headersChanged := false
	if response.ETag != "" && response.ETag != blog.ETag() {
		headersChanged = true
		blog.SetETag(response.ETag)
	}

	if response.LastModified != "" && response.LastModified != blog.LastModified() {
		headersChanged = true
		blog.SetLastModified(response.LastModified)
	}

	// Update the blog's cache headers if changed.
	if headersChanged {
		return s.repo.Blog().Update(blog)
	}

	return nil
}

func (s *SyncService) Run(ctx context.Context) error {
	// perform an initial sync at service startup
	err := s.SyncAllBlogs()
	if err != nil {
		slog.Error("error syncing blogs",
			"error", err.Error(),
		)
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
				slog.Error("error syncing blogs",
					"error", err.Error(),
				)
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

	blogs, err := s.repo.Blog().ListAll()
	if err != nil {
		return err
	}

	// Be sure to only sync blogs that are ready.
	now := timeutil.Now()
	syncableBlogs := FilterSyncableBlogs(blogs, now)

	// Update the syncedAt time for each blog before syncing.
	for _, blog := range syncableBlogs {
		blog.SetSyncedAt(now)
		err = s.repo.Blog().Update(blog)
		if err != nil {
			return err
		}
	}

	// use a weighted semaphore to limit concurrency
	sem := semaphore.NewWeighted(SyncConcurrency)

	// sync all blogs in parallel (up to SyncConcurrency at once)
	for _, blog := range syncableBlogs {
		sem.Acquire(context.Background(), 1)

		go func(blog *model.Blog) {
			defer sem.Release(1)

			slog.Info("syncing blog", "title", blog.Title(), "id", blog.ID())
			_, err := s.SyncBlog(blog.FeedURL())
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

// Sync a new or existing Blog based on the provided feed URL.
func (s *SyncService) SyncBlog(feedURL string) (*model.Blog, error) {
	blog, err := s.repo.Blog().ReadByFeedURL(feedURL)
	if err != nil {
		if !errors.Is(err, postgres.ErrNotFound) {
			return nil, err
		}

		return s.syncNewBlog(feedURL)
	}

	return s.syncExistingBlog(blog)
}

func (s *SyncService) syncNewBlog(feedURL string) (*model.Blog, error) {
	req := fetch.FetchFeedRequest{
		URL: feedURL,
	}
	resp, err := s.feedFetcher.FetchFeed(req)
	if err != nil {
		return nil, err
	}

	// no feed from a new blog is an error
	if resp.Feed == "" {
		return nil, fetch.ErrUnreachableFeed
	}

	feedBlog, err := feed.Parse(feedURL, resp.Feed)
	if err != nil {
		return nil, err
	}

	blog, err := model.NewBlog(
		feedBlog.FeedURL,
		feedBlog.SiteURL,
		feedBlog.Title,
		resp.ETag,
		resp.LastModified,
		timeutil.Now(),
	)
	if err != nil {
		return nil, err
	}

	err = s.repo.Blog().Create(blog)
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
	err := s.repo.Blog().Update(blog)
	if err != nil {
		return nil, err
	}

	req := fetch.FetchFeedRequest{
		URL:          blog.FeedURL(),
		ETag:         blog.ETag(),
		LastModified: blog.LastModified(),
	}
	resp, err := s.feedFetcher.FetchFeed(req)
	if err != nil {
		return nil, err
	}

	err = s.updateCacheHeaders(blog, resp)
	if err != nil {
		return nil, err
	}

	if resp.Feed == "" {
		slog.Info("skipping blog (no feed content)", "title", blog.Title(), "id", blog.ID())
		return blog, nil
	}

	feedBlog, err := feed.Parse(blog.FeedURL(), resp.Feed)
	if err != nil {
		return nil, err
	}

	for _, feedPost := range feedBlog.Posts {
		_, err = s.syncPost(blog, feedPost)
		if err != nil {
			slog.Warn("failed to sync post", "url", feedPost.URL, "error", err.Error())
		}
	}

	return blog, nil
}

func (s *SyncService) syncPost(blog *model.Blog, feedPost feed.Post) (*model.Post, error) {
	post, err := s.repo.Post().ReadByURL(feedPost.URL)
	if err != nil {
		if !errors.Is(err, postgres.ErrNotFound) {
			return nil, err
		}

		// If the post doesn't exist, create it and fall through.
		post, err = model.NewPost(
			blog,
			feedPost.URL,
			feedPost.Title,
			feedPost.Content,
			feedPost.PublishedAt,
		)
		if err != nil {
			return nil, err
		}

		err = s.repo.Post().Create(post)
		if err != nil {
			return nil, err
		}
	}

	// Update the post's content (if not present but available from the feed).
	if post.Content() == "" && feedPost.Content != "" {
		post.SetContent(feedPost.Content)
		err = s.repo.Post().Update(post)
		if err != nil {
			return nil, err
		}
	}

	// If we still don't have content, try to fetch it manually.
	if post.Content() == "" {
		req := fetch.FetchPageRequest{
			URL: post.URL(),
		}
		resp, err := s.pageFetcher.FetchPage(req)
		if err != nil {
			return post, nil
		}

		post.SetContent(resp.Content)
		err = s.repo.Post().Update(post)
		if err != nil {
			return nil, err
		}
	}

	return post, nil
}
