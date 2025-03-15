package service_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/feed"
	feedMock "github.com/theandrew168/bloggulus/backend/feed/mock"
	"github.com/theandrew168/bloggulus/backend/fetch"
	fetchMock "github.com/theandrew168/bloggulus/backend/fetch/mock"
	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

func TestFilterSyncableBlogs(t *testing.T) {
	t.Parallel()

	now := timeutil.Now()

	pastBlog := test.NewBlog(t)
	pastBlog.SetSyncedAt(now.Add(-model.SyncCooldown).Add(-1 * time.Minute))
	presentBlog := test.NewBlog(t)
	presentBlog.SetSyncedAt(now)
	futureBlog := test.NewBlog(t)
	futureBlog.SetSyncedAt(now.Add(1 * time.Hour))

	blogs := []*model.Blog{pastBlog, presentBlog, futureBlog}

	syncableBlogs := service.FilterSyncableBlogs(blogs, now)
	test.AssertEqual(t, len(syncableBlogs), 1)

	var syncableBlogIDs []uuid.UUID
	for _, blog := range syncableBlogs {
		syncableBlogIDs = append(syncableBlogIDs, blog.ID())
	}

	test.AssertSliceContains(t, syncableBlogIDs, pastBlog.ID())
}

func TestUpdateCacheHeaders(t *testing.T) {
	t.Parallel()

	blog := test.NewBlog(t)
	resp := fetch.FetchFeedResponse{
		ETag:         "foo",
		LastModified: "bar",
	}

	changed := service.UpdateCacheHeaders(blog, resp)
	test.AssertEqual(t, changed, true)
	test.AssertEqual(t, blog.ETag(), "foo")
	test.AssertEqual(t, blog.LastModified(), "bar")
}

func TestUpdateCacheHeadersDoesNotClear(t *testing.T) {
	t.Parallel()

	blog := test.NewBlog(t)
	resp := fetch.FetchFeedResponse{
		ETag:         "",
		LastModified: "",
	}

	changed := service.UpdateCacheHeaders(blog, resp)
	test.AssertEqual(t, changed, false)
	test.AssertEqual(t, blog.ETag(), blog.ETag())
	test.AssertEqual(t, blog.LastModified(), blog.LastModified())
}

func TestComparePosts(t *testing.T) {
	t.Parallel()

	blog := test.NewBlog(t)

	knownPost := test.NewPost(t, blog)
	knownPosts := []*model.Post{
		knownPost,
	}

	newPost := feed.Post{
		URL:         test.RandomURL(20),
		Title:       test.RandomString(20),
		Content:     test.RandomString(200),
		PublishedAt: test.RandomTime(),
	}
	updatedPost := feed.Post{
		URL:         knownPost.URL(),
		Title:       test.RandomString(20),
		Content:     test.RandomString(200),
		PublishedAt: test.RandomTime(),
	}
	feedPosts := []feed.Post{
		newPost,
		updatedPost,
	}

	result, err := service.ComparePosts(blog, knownPosts, feedPosts)
	test.AssertNilError(t, err)

	// Verify that one new post should be created.
	test.AssertEqual(t, len(result.PostsToCreate), 1)
	test.AssertEqual(t, result.PostsToCreate[0].URL(), newPost.URL)
	test.AssertEqual(t, result.PostsToCreate[0].Title(), newPost.Title)
	test.AssertEqual(t, result.PostsToCreate[0].Content(), newPost.Content)
	test.AssertEqual(t, result.PostsToCreate[0].PublishedAt(), newPost.PublishedAt)

	// Verify that one existing post should be updated (URL should stay the same).
	test.AssertEqual(t, len(result.PostsToUpdate), 1)
	test.AssertEqual(t, result.PostsToUpdate[0].URL(), knownPost.URL())
	test.AssertEqual(t, result.PostsToUpdate[0].URL(), updatedPost.URL)
	test.AssertEqual(t, result.PostsToUpdate[0].Title(), updatedPost.Title)
	test.AssertEqual(t, result.PostsToUpdate[0].Content(), updatedPost.Content)
	test.AssertEqual(t, result.PostsToUpdate[0].PublishedAt(), updatedPost.PublishedAt)
}

func TestNewBlog(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	feedPost := feed.Post{
		URL:         test.RandomURL(20),
		Title:       test.RandomString(20),
		Content:     test.RandomString(200),
		PublishedAt: time.Now(),
	}
	feedBlog := feed.Blog{
		Title:   test.RandomString(20),
		SiteURL: test.RandomString(20),
		FeedURL: test.RandomString(20),
		Posts:   []feed.Post{feedPost},
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	feeds := map[string]fetch.FetchFeedResponse{
		feedBlog.FeedURL: {Feed: atomFeed},
	}
	feedFetcher := fetchMock.NewFeedFetcher(feeds)

	syncService := service.NewSyncService(repo, feedFetcher)

	// sync a new blog
	blog, err := syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// verify blog data
	test.AssertNilError(t, err)
	test.AssertEqual(t, blog.Title(), feedBlog.Title)
	test.AssertEqual(t, blog.SiteURL(), feedBlog.SiteURL)
	test.AssertEqual(t, blog.FeedURL(), feedBlog.FeedURL)

	// fetch posts and verify count
	posts, err := repo.Post().ListByBlog(blog)
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(posts), 1)

	// verify post data
	post := posts[0]
	test.AssertEqual(t, post.URL(), feedPost.URL)
	test.AssertEqual(t, post.Title(), feedPost.Title)
	test.AssertEqual(t, post.Content(), feedPost.Content)
}

func TestExistingBlog(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	feedBlog := feed.Blog{
		Title:   test.RandomString(20),
		SiteURL: test.RandomURL(20),
		FeedURL: test.RandomURL(20),
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	feeds := map[string]fetch.FetchFeedResponse{
		feedBlog.FeedURL: {Feed: atomFeed},
	}
	feedFetcher := fetchMock.NewFeedFetcher(feeds)

	syncService := service.NewSyncService(repo, feedFetcher)

	// sync a new blog
	blog, err := syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// verify blog data
	test.AssertNilError(t, err)
	test.AssertEqual(t, blog.Title(), feedBlog.Title)
	test.AssertEqual(t, blog.SiteURL(), feedBlog.SiteURL)
	test.AssertEqual(t, blog.FeedURL(), feedBlog.FeedURL)

	// fetch posts and verify count (should be none)
	posts, err := repo.Post().ListByBlog(blog)
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(posts), 0)

	// add a post to the feed blog
	feedPost := feed.Post{
		URL:         test.RandomURL(20),
		Title:       test.RandomString(20),
		Content:     test.RandomString(200),
		PublishedAt: time.Now(),
	}
	feedBlog.Posts = append(feedBlog.Posts, feedPost)

	// regenerate the feed
	atomFeed, err = feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	feeds[feedBlog.FeedURL] = fetch.FetchFeedResponse{Feed: atomFeed}

	// sync the blog again
	_, err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// fetch posts and verify count
	posts, err = repo.Post().ListByBlog(blog)
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(posts), 1)

	// verify post data
	post := posts[0]
	test.AssertEqual(t, post.URL(), feedPost.URL)
	test.AssertEqual(t, post.Title(), feedPost.Title)
	test.AssertEqual(t, post.Content(), feedPost.Content)
}

func TestUnreachableFeed(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	feedURL := test.RandomURL(20)

	feeds := map[string]fetch.FetchFeedResponse{}
	feedFetcher := fetchMock.NewFeedFetcher(feeds)

	syncService := service.NewSyncService(repo, feedFetcher)

	_, err := syncService.SyncBlog(feedURL)
	test.AssertErrorIs(t, err, fetch.ErrUnreachableFeed)
}

func TestSyncCooldown(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	feedBlog := feed.Blog{
		Title:   test.RandomString(20),
		SiteURL: test.RandomURL(20),
		FeedURL: test.RandomURL(20),
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	feeds := map[string]fetch.FetchFeedResponse{
		feedBlog.FeedURL: {Feed: atomFeed},
	}
	feedFetcher := fetchMock.NewFeedFetcher(feeds)

	syncService := service.NewSyncService(repo, feedFetcher)

	// add a blog (sync now)
	blog, err := syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// capture the blog's current syncedAt time
	syncedAt := blog.SyncedAt()

	// sync all blogs
	err = syncService.SyncAllBlogs()
	test.AssertNilError(t, err)

	// refetch the blog's data
	blog, err = repo.Blog().ReadByFeedURL(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// syncedAt should not have changed
	test.AssertEqual(t, blog.SyncedAt(), syncedAt)
}

func TestUpdatePostContent(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	feedPost := feed.Post{
		URL:         test.RandomURL(20),
		Title:       test.RandomString(20),
		PublishedAt: time.Now(),
	}
	feedBlog := feed.Blog{
		Title:   test.RandomString(20),
		SiteURL: test.RandomURL(20),
		FeedURL: test.RandomURL(20),
		Posts:   []feed.Post{feedPost},
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	feeds := map[string]fetch.FetchFeedResponse{
		feedBlog.FeedURL: {Feed: atomFeed},
	}
	feedFetcher := fetchMock.NewFeedFetcher(feeds)

	syncService := service.NewSyncService(repo, feedFetcher)

	// sync a new blog
	blog, err := syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// fetch posts and verify count
	posts, err := repo.Post().ListByBlog(blog)
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(posts), 1)

	// verify post data (should have no content)
	post := posts[0]
	test.AssertEqual(t, post.Content(), "")

	// update the post with some content
	content := "content about foo"
	feedBlog.Posts[0].Content = content

	// regenerate the feed
	atomFeed, err = feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	feeds[feedBlog.FeedURL] = fetch.FetchFeedResponse{Feed: atomFeed}

	// sync the blog again
	_, err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// refetch posts and verify count
	posts, err = repo.Post().ListByBlog(blog)
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(posts), 1)

	// verify post data (should have content now)
	post = posts[0]
	test.AssertEqual(t, post.Content(), content)
}

// Doesn't wipe out existing cache headers if none are returned.
func TestCacheHeaderOverwrite(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	feedBlog := feed.Blog{
		Title:   test.RandomString(20),
		SiteURL: test.RandomURL(20),
		FeedURL: test.RandomURL(20),
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	feeds := map[string]fetch.FetchFeedResponse{
		feedBlog.FeedURL: {Feed: atomFeed},
	}
	feedFetcher := fetchMock.NewFeedFetcher(feeds)

	syncService := service.NewSyncService(repo, feedFetcher)

	// sync a new blog
	blog, err := syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// update the blog's ETag and LastModified to something non-empty
	blog.SetETag("foo")
	blog.SetLastModified("bar")
	err = repo.Blog().Update(blog)
	test.AssertNilError(t, err)

	// sync the blog again (will see empty ETag and LastModified values)
	_, err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// refetch the blog
	blog, err = repo.Blog().ReadByFeedURL(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// verify that the existing ETag and LastModified values haven't been wiped out
	test.AssertEqual(t, blog.ETag(), "foo")
	test.AssertEqual(t, blog.LastModified(), "bar")
}

// Update cache header values even if no new content is available.
func TestCacheHeaderUpdate(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	feedBlog := feed.Blog{
		Title:   test.RandomString(20),
		SiteURL: test.RandomURL(20),
		FeedURL: test.RandomURL(20),
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	feeds := map[string]fetch.FetchFeedResponse{
		feedBlog.FeedURL: {
			Feed:         atomFeed,
			ETag:         "etag",
			LastModified: "lastModified",
		},
	}
	feedFetcher := fetchMock.NewFeedFetcher(feeds)

	syncService := service.NewSyncService(repo, feedFetcher)

	// sync a new blog
	blog, err := syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)
	test.AssertEqual(t, blog.ETag(), "etag")
	test.AssertEqual(t, blog.LastModified(), "lastModified")

	// Update the feed to return new cache value but no data
	feeds = map[string]fetch.FetchFeedResponse{
		feedBlog.FeedURL: {
			ETag:         "other etag",
			LastModified: "other lastModified",
		},
	}
	feedFetcher = fetchMock.NewFeedFetcher(feeds)

	syncService = service.NewSyncService(repo, feedFetcher)

	// sync the blog again (will see new ETag and LastModified values)
	_, err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// refetch the blog
	blog, err = repo.Blog().ReadByFeedURL(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// verify that the ETag and LastModified values got updated
	test.AssertEqual(t, blog.ETag(), "other etag")
	test.AssertEqual(t, blog.LastModified(), "other lastModified")
}
