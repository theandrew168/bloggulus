package service_test

import (
	"testing"
	"time"

	"github.com/theandrew168/bloggulus/backend/domain/admin/feed"
	feedMock "github.com/theandrew168/bloggulus/backend/domain/admin/feed/mock"
	"github.com/theandrew168/bloggulus/backend/domain/admin/fetch"
	fetchMock "github.com/theandrew168/bloggulus/backend/domain/admin/fetch/mock"
	"github.com/theandrew168/bloggulus/backend/domain/admin/service"
	storageMock "github.com/theandrew168/bloggulus/backend/domain/admin/storage/mock"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestNewBlog(t *testing.T) {
	t.Parallel()

	feedPost := feed.Post{
		URL:         "https://example.com/foo",
		Title:       "Foo",
		Content:     "content about foo",
		PublishedAt: time.Now(),
	}
	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []feed.Post{feedPost},
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	store := storageMock.NewStorage()

	feeds := map[string]string{
		feedBlog.FeedURL: atomFeed,
	}
	feedFetcher := fetchMock.NewFeedFetcher(feeds)

	pages := map[string]string{}
	pageFetcher := fetchMock.NewPageFetcher(pages)

	syncService := service.NewSyncService(store, feedFetcher, pageFetcher)

	// sync a new blog
	err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// fetch and verify blog data
	blog, err := store.Blog().ReadByFeedURL(feedBlog.FeedURL)
	test.AssertNilError(t, err)
	test.AssertEqual(t, blog.Title(), feedBlog.Title)
	test.AssertEqual(t, blog.SiteURL(), feedBlog.SiteURL)
	test.AssertEqual(t, blog.FeedURL(), feedBlog.FeedURL)

	// fetch posts and verify count
	posts, err := store.Post().ListByBlog(blog, 20, 0)
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

	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	store := storageMock.NewStorage()

	feeds := map[string]string{
		feedBlog.FeedURL: atomFeed,
	}
	feedFetcher := fetchMock.NewFeedFetcher(feeds)

	pages := map[string]string{}
	pageFetcher := fetchMock.NewPageFetcher(pages)

	syncService := service.NewSyncService(store, feedFetcher, pageFetcher)

	// sync a new blog
	err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// fetch and verify blog data
	blog, err := store.Blog().ReadByFeedURL(feedBlog.FeedURL)
	test.AssertNilError(t, err)
	test.AssertEqual(t, blog.Title(), feedBlog.Title)
	test.AssertEqual(t, blog.SiteURL(), feedBlog.SiteURL)
	test.AssertEqual(t, blog.FeedURL(), feedBlog.FeedURL)

	// fetch posts and verify count (should be none)
	posts, err := store.Post().ListByBlog(blog, 20, 0)
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(posts), 0)

	// add a post to the feed blog
	feedPost := feed.Post{
		URL:         "https://example.com/foo",
		Title:       "Foo",
		Content:     "content about foo",
		PublishedAt: time.Now(),
	}
	feedBlog.Posts = append(feedBlog.Posts, feedPost)

	// regenerate the feed
	atomFeed, err = feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	feeds[feedBlog.FeedURL] = atomFeed

	// sync the blog again
	err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// fetch posts and verify count
	posts, err = store.Post().ListByBlog(blog, 20, 0)
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

	feedURL := "https://example.com/atom.xml"

	store := storageMock.NewStorage()

	feeds := map[string]string{}
	feedFetcher := fetchMock.NewFeedFetcher(feeds)

	pages := map[string]string{}
	pageFetcher := fetchMock.NewPageFetcher(pages)

	syncService := service.NewSyncService(store, feedFetcher, pageFetcher)

	err := syncService.SyncBlog(feedURL)
	test.AssertErrorIs(t, err, fetch.ErrUnreachableFeed)
}

func TestNoNewFeedContent(t *testing.T) {
	t.Parallel()

	feedURL := "https://example.com/atom.xml"

	store := storageMock.NewStorage()

	feeds := map[string]string{
		feedURL: "",
	}
	feedFetcher := fetchMock.NewFeedFetcher(feeds)

	pages := map[string]string{}
	pageFetcher := fetchMock.NewPageFetcher(pages)

	syncService := service.NewSyncService(store, feedFetcher, pageFetcher)

	err := syncService.SyncBlog(feedURL)
	test.AssertErrorIs(t, err, fetch.ErrNoNewFeedContent)
}

func TestSyncOncePerHour(t *testing.T) {
	t.Parallel()

	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	store := storageMock.NewStorage()

	feeds := map[string]string{
		feedBlog.FeedURL: atomFeed,
	}
	feedFetcher := fetchMock.NewFeedFetcher(feeds)

	pages := map[string]string{}
	pageFetcher := fetchMock.NewPageFetcher(pages)

	syncService := service.NewSyncService(store, feedFetcher, pageFetcher)

	// add a blog (sync now)
	err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// fetch the synced blog
	blog, err := store.Blog().ReadByFeedURL(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// capture its current syncedAt time
	syncedAt := blog.SyncedAt()

	// sync all blogs
	err = syncService.SyncAllBlogs()
	test.AssertNilError(t, err)

	// refetch the blog's data
	blog, err = store.Blog().ReadByFeedURL(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// syncedAt should not have changed
	test.AssertEqual(t, blog.SyncedAt(), syncedAt)
}

func TestUpdatePostContent(t *testing.T) {
	t.Parallel()

	feedPost := feed.Post{
		URL:         "https://example.com/foo",
		Title:       "Foo",
		PublishedAt: time.Now(),
	}
	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []feed.Post{feedPost},
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	store := storageMock.NewStorage()

	feeds := map[string]string{
		feedBlog.FeedURL: atomFeed,
	}
	feedFetcher := fetchMock.NewFeedFetcher(feeds)

	pages := map[string]string{}
	pageFetcher := fetchMock.NewPageFetcher(pages)

	syncService := service.NewSyncService(store, feedFetcher, pageFetcher)

	// sync a new blog
	err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// fetch the synced blog
	blog, err := store.Blog().ReadByFeedURL(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// fetch posts and verify count
	posts, err := store.Post().ListByBlog(blog, 20, 0)
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

	feeds[feedBlog.FeedURL] = atomFeed

	// sync the blog again
	err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// refetch posts and verify count
	posts, err = store.Post().ListByBlog(blog, 20, 0)
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(posts), 1)

	// verify post data (should have content now)
	post = posts[0]
	test.AssertEqual(t, post.Content(), content)
}

func TestCacheHeaderOverwrite(t *testing.T) {
	t.Parallel()

	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	store := storageMock.NewStorage()

	feeds := map[string]string{
		feedBlog.FeedURL: atomFeed,
	}
	feedFetcher := fetchMock.NewFeedFetcher(feeds)

	pages := map[string]string{}
	pageFetcher := fetchMock.NewPageFetcher(pages)

	syncService := service.NewSyncService(store, feedFetcher, pageFetcher)

	// sync a new blog
	err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// fetch the synced blog
	blog, err := store.Blog().ReadByFeedURL(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// update the blog's ETag and LastModified to something non-empty
	blog.SetETag("foo")
	blog.SetLastModified("bar")
	err = store.Blog().Update(blog)
	test.AssertNilError(t, err)

	// sync the block again (will see empty ETag and LastModified values)
	err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// refetch the blog
	blog, err = store.Blog().ReadByFeedURL(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	// verify that the existing ETag and LastModified values haven't been wiped out
	test.AssertEqual(t, blog.ETag(), "foo")
	test.AssertEqual(t, blog.LastModified(), "bar")
}
