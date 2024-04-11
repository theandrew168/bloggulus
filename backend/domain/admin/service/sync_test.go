package service_test

import (
	"testing"
	"time"

	"github.com/theandrew168/bloggulus/backend/domain/admin/feed"
	feedMock "github.com/theandrew168/bloggulus/backend/domain/admin/feed/mock"
	fetchMock "github.com/theandrew168/bloggulus/backend/domain/admin/fetch/mock"
	"github.com/theandrew168/bloggulus/backend/domain/admin/service"
	storageMock "github.com/theandrew168/bloggulus/backend/domain/admin/storage/mock"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestNewBlog(t *testing.T) {
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

	err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	blog, err := store.Blog().ReadByFeedURL(feedBlog.FeedURL)
	test.AssertNilError(t, err)
	test.AssertEqual(t, blog.Title(), feedBlog.Title)
	test.AssertEqual(t, blog.SiteURL(), feedBlog.SiteURL)
	test.AssertEqual(t, blog.FeedURL(), feedBlog.FeedURL)

	posts, err := store.Post().ListByBlog(blog, 20, 0)
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(posts), 1)

	post := posts[0]
	test.AssertEqual(t, post.URL(), feedPost.URL)
	test.AssertEqual(t, post.Title(), feedPost.Title)
	test.AssertEqual(t, post.Content(), feedPost.Content)
}

func TestExistingBlog(t *testing.T) {
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

	err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	blog, err := store.Blog().ReadByFeedURL(feedBlog.FeedURL)
	test.AssertNilError(t, err)
	test.AssertEqual(t, blog.Title(), feedBlog.Title)
	test.AssertEqual(t, blog.SiteURL(), feedBlog.SiteURL)
	test.AssertEqual(t, blog.FeedURL(), feedBlog.FeedURL)

	posts, err := store.Post().ListByBlog(blog, 20, 0)
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(posts), 0)

	feedPost := feed.Post{
		URL:         "https://example.com/foo",
		Title:       "Foo",
		Content:     "content about foo",
		PublishedAt: time.Now(),
	}
	feedBlog.Posts = append(feedBlog.Posts, feedPost)

	atomFeed, err = feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	feeds[feedBlog.FeedURL] = atomFeed

	err = syncService.SyncBlog(feedBlog.FeedURL)
	test.AssertNilError(t, err)

	posts, err = store.Post().ListByBlog(blog, 20, 0)
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(posts), 1)

	post := posts[0]
	test.AssertEqual(t, post.URL(), feedPost.URL)
	test.AssertEqual(t, post.Title(), feedPost.Title)
	test.AssertEqual(t, post.Content(), feedPost.Content)
}

func TestSkipEmptyFeed(t *testing.T) {

}

func TestSyncOncePerHour(t *testing.T) {

}

func TestUpdatePostContent(t *testing.T) {

}

func TestSkipPostsWithoutLink(t *testing.T) {

}

func TestSkipPostsWithoutTitle(t *testing.T) {

}

func TestCacheHeaderOverwrite(t *testing.T) {

}
