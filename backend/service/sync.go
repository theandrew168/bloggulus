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
	SyncInterval = 30 * time.Minute

	// How many blogs to sync at once.
	SyncConcurrency = 8
)

// FilterSyncableBlogs takes a list of blogs and returns only those that are ready to be synced.
func FilterSyncableBlogs(blogs []*model.Blog, now time.Time) []*model.Blog {
	var syncableBlogs []*model.Blog
	for _, blog := range blogs {
		if blog.CanBeSynced(now) {
			syncableBlogs = append(syncableBlogs, blog)
		}
	}
	return syncableBlogs
}

// UpdateCacheHeaders updates the ETag and Last-Modified headers for a blog if they have changed.
func UpdateCacheHeaders(blog *model.Blog, response fetch.FetchFeedResponse) bool {
	headersChanged := false
	if response.ETag != "" && response.ETag != blog.ETag() {
		headersChanged = true
		blog.SetETag(response.ETag)
	}

	if response.LastModified != "" && response.LastModified != blog.LastModified() {
		headersChanged = true
		blog.SetLastModified(response.LastModified)
	}

	return headersChanged
}

type ComparePostsResult struct {
	PostsToCreate []*model.Post
	PostsToUpdate []*model.Post
}

// ComparePosts compares a list of known posts to a list of feed posts and returns
// a list of posts to create and a list of posts to update.
func ComparePosts(blog *model.Blog, knownPosts []*model.Post, feedPosts []feed.Post) (ComparePostsResult, error) {
	// Create a map of URLs to posts for quick lookups.
	knownPostsByURL := make(map[string]*model.Post)
	for _, post := range knownPosts {
		knownPostsByURL[post.URL()] = post
	}

	var postsToCreate []*model.Post
	var postsToUpdate []*model.Post

	// Compare each post in the feed to the posts in the database.
	for _, feedPost := range feedPosts {
		knownPost, ok := knownPostsByURL[feedPost.URL]
		if !ok {
			// The post is new so we need to create it.
			postToCreate, err := model.NewPost(
				blog,
				feedPost.URL,
				feedPost.Title,
				feedPost.Content,
				feedPost.PublishedAt,
			)
			if err != nil {
				return ComparePostsResult{}, err
			}

			postsToCreate = append(postsToCreate, postToCreate)
		} else {
			// The post already exists but we might need to update it.
			knownPostShouldBeUpdated := false

			// Check if the post's title has changed.
			if feedPost.Title != "" && feedPost.Title != knownPost.Title() {
				knownPost.SetTitle(feedPost.Title)
				knownPostShouldBeUpdated = true
			}

			// Check if the post's content has changed.
			if feedPost.Content != "" && feedPost.Content != knownPost.Content() {
				knownPost.SetContent(feedPost.Content)
				knownPostShouldBeUpdated = true
			}

			// Check if the post's publishedAt date has changed.
			if feedPost.PublishedAt != knownPost.PublishedAt() {
				knownPost.SetPublishedAt(feedPost.PublishedAt)
				knownPostShouldBeUpdated = true
			}

			// If any post data has changed, add it to the list of posts to update.
			if knownPostShouldBeUpdated {
				postsToUpdate = append(postsToUpdate, knownPost)
			}
		}
	}

	result := ComparePostsResult{
		PostsToCreate: postsToCreate,
		PostsToUpdate: postsToUpdate,
	}
	return result, nil
}

func ParallelMap[T any](concurrency int, items []T, fn func(T)) {
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

type SyncService struct {
	mu          sync.Mutex
	repo        *repository.Repository
	feedFetcher fetch.FeedFetcher
}

func NewSyncService(repo *repository.Repository, feedFetcher fetch.FeedFetcher) *SyncService {
	s := SyncService{
		repo:        repo,
		feedFetcher: feedFetcher,
	}
	return &s
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

	// Update the syncedAt time for each syncable blog before syncing.
	for _, blog := range syncableBlogs {
		blog.SetSyncedAt(now)
		err = s.repo.Blog().Update(blog)
		if err != nil {
			return err
		}
	}

	ParallelMap(SyncConcurrency, syncableBlogs, func(blog *model.Blog) {
		slog.Info("syncing blog", "title", blog.Title(), "id", blog.ID())
		_, err := s.SyncBlog(blog.FeedURL())
		if err != nil {
			slog.Warn(err.Error(), "title", blog.Title(), "id", blog.ID())
		}
	})

	return nil
}

// Sync a new or existing Blog based on the provided feed URL.
func (s *SyncService) SyncBlog(feedURL string) (*model.Blog, error) {
	blog, err := s.repo.Blog().ReadByFeedURL(feedURL)
	if err != nil {
		if !errors.Is(err, postgres.ErrNotFound) {
			return nil, err
		}

		// An ErrNotFound is acceptable (and expected) here. The only difference
		// is that we won't be able to include the ETag and Last-Modified headers
		// in the request. This is fine for new blogs (an unconditional fetch).
		return s.syncNewBlog(feedURL)
	}

	return s.syncExistingBlog(blog)
}

func (s *SyncService) syncNewBlog(feedURL string) (*model.Blog, error) {
	// Make an unconditional fetch for the blog's feed.
	req := fetch.FetchFeedRequest{
		URL: feedURL,
	}
	resp, err := s.feedFetcher.FetchFeed(req)
	if err != nil {
		return nil, err
	}

	// No feed data from a new blog is an error.
	if resp.Feed == "" {
		return nil, fetch.ErrUnreachableFeed
	}

	feedBlog, err := feed.Parse(feedURL, resp.Feed)
	if err != nil {
		return nil, err
	}

	// Create a new blog based on the feed data.
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

	err = s.syncPosts(blog, feedBlog.Posts)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *SyncService) syncExistingBlog(blog *model.Blog) (*model.Blog, error) {
	// Make a conditional fetch for the blog's feed.
	req := fetch.FetchFeedRequest{
		URL:          blog.FeedURL(),
		ETag:         blog.ETag(),
		LastModified: blog.LastModified(),
	}
	resp, err := s.feedFetcher.FetchFeed(req)
	if err != nil {
		return nil, err
	}

	// Update the blog's cache headers if they have changed.
	headersChanged := UpdateCacheHeaders(blog, resp)
	if headersChanged {
		err = s.repo.Blog().Update(blog)
		if err != nil {
			return nil, err
		}
	}

	if resp.Feed == "" {
		slog.Info("skipping blog (no feed content)", "title", blog.Title(), "id", blog.ID())
		return blog, nil
	}

	feedBlog, err := feed.Parse(blog.FeedURL(), resp.Feed)
	if err != nil {
		return nil, err
	}

	err = s.syncPosts(blog, feedBlog.Posts)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *SyncService) syncPosts(blog *model.Blog, feedPosts []feed.Post) error {
	// List all known posts for the current blog.
	knownPosts, err := s.repo.Post().ListByBlog(blog)
	if err != nil {
		return err
	}

	// Compare the known posts to the feed posts.
	result, err := ComparePosts(blog, knownPosts, feedPosts)
	if err != nil {
		return err
	}

	// Create any posts that are new.
	for _, post := range result.PostsToCreate {
		err = s.repo.Post().Create(post)
		if err != nil {
			slog.Warn("failed to create post", "url", post.URL(), "error", err.Error())
		}
	}

	// Update any posts that have changed.
	for _, post := range result.PostsToUpdate {
		err = s.repo.Post().Update(post)
		if err != nil {
			slog.Warn("failed to update post", "url", post.URL(), "error", err.Error())
		}
	}

	return nil
}
