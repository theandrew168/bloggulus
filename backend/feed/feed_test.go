package feed_test

import (
	"testing"
	"time"

	"github.com/theandrew168/bloggulus/backend/feed"
	feedMock "github.com/theandrew168/bloggulus/backend/feed/mock"
	fetchMock "github.com/theandrew168/bloggulus/backend/fetch/mock"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestParse(t *testing.T) {
	t.Parallel()

	feedPostFoo := feed.Post{
		URL:         "https://example.com/foo",
		Title:       "Foo",
		Content:     "content about foo",
		PublishedAt: time.Now(),
	}
	feedPostBar := feed.Post{
		URL:         "https://example.com/bar",
		Title:       "Bar",
		Content:     "content about bar",
		PublishedAt: time.Now(),
	}
	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []feed.Post{feedPostFoo, feedPostBar},
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	parsedBlog, err := feed.Parse("https://example.com/atom.xml", atomFeed)
	test.AssertNilError(t, err)
	test.AssertEqual(t, parsedBlog.Title, feedBlog.Title)
	test.AssertEqual(t, parsedBlog.SiteURL, feedBlog.SiteURL)
	test.AssertEqual(t, parsedBlog.FeedURL, feedBlog.FeedURL)
	test.AssertEqual(t, len(parsedBlog.Posts), len(feedBlog.Posts))

	postsByURL := map[string]feed.Post{
		feedPostFoo.URL: feedPostFoo,
		feedPostBar.URL: feedPostBar,
	}
	for _, parsedPost := range parsedBlog.Posts {
		post, ok := postsByURL[parsedPost.URL]
		if !ok {
			t.Errorf("invalid post URL: %s", parsedPost.URL)
			continue
		}

		test.AssertEqual(t, parsedPost.URL, post.URL)
		test.AssertEqual(t, parsedPost.Title, post.Title)
		test.AssertEqual(t, parsedPost.Content, post.Content)
	}
}

func TestParseMissingURL(t *testing.T) {
	t.Parallel()

	feedPostFoo := feed.Post{
		Title:       "Foo",
		Content:     "content about foo",
		PublishedAt: time.Now(),
	}
	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []feed.Post{feedPostFoo},
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	parsedBlog, err := feed.Parse("https://example.com/atom.xml", atomFeed)
	test.AssertNilError(t, err)

	test.AssertEqual(t, len(parsedBlog.Posts), 0)
}

func TestParseMissingTitle(t *testing.T) {
	t.Parallel()

	feedPostFoo := feed.Post{
		URL:         "https://example.com/foo",
		Content:     "content about foo",
		PublishedAt: time.Now(),
	}
	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []feed.Post{feedPostFoo},
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	parsedBlog, err := feed.Parse("https://example.com/atom.xml", atomFeed)
	test.AssertNilError(t, err)

	test.AssertEqual(t, len(parsedBlog.Posts), 0)
}

func TestParseMissingDomain(t *testing.T) {
	t.Parallel()

	feedPostFoo := feed.Post{
		URL:         "/foo",
		Title:       "Foo",
		Content:     "content about foo",
		PublishedAt: time.Now(),
	}
	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []feed.Post{feedPostFoo},
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	parsedBlog, err := feed.Parse("https://example.com/atom.xml", atomFeed)
	test.AssertNilError(t, err)

	for _, parsedPost := range parsedBlog.Posts {
		test.AssertEqual(t, parsedPost.URL, feedBlog.SiteURL+feedPostFoo.URL)
	}
}

func TestParseMissingScheme(t *testing.T) {
	t.Parallel()

	feedPostFoo := feed.Post{
		URL:         "example.com/foo",
		Title:       "Foo",
		Content:     "content about foo",
		PublishedAt: time.Now(),
	}
	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []feed.Post{feedPostFoo},
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	parsedBlog, err := feed.Parse("https://example.com/atom.xml", atomFeed)
	test.AssertNilError(t, err)

	for _, parsedPost := range parsedBlog.Posts {
		test.AssertEqual(t, parsedPost.URL, "https://"+feedPostFoo.URL)
	}
}

func TestParsePublishedAtUTC(t *testing.T) {
	t.Parallel()

	publishedAt, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
	test.AssertNilError(t, err)

	feedPostFoo := feed.Post{
		URL:         "https://example.com/foo",
		Title:       "Foo",
		Content:     "content about foo",
		PublishedAt: publishedAt,
	}
	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []feed.Post{feedPostFoo},
	}

	atomFeed, err := feedMock.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	parsedBlog, err := feed.Parse("https://example.com/atom.xml", atomFeed)
	test.AssertNilError(t, err)

	for _, parsedPost := range parsedBlog.Posts {
		test.AssertEqual(t, parsedPost.PublishedAt, publishedAt.UTC().Round(time.Microsecond))
	}
}

func TestHydrate(t *testing.T) {
	t.Parallel()

	feedPostFoo := feed.Post{
		URL:         "https://example.com/foo",
		Title:       "Foo",
		PublishedAt: time.Now(),
	}
	feedPostBar := feed.Post{
		URL:         "https://example.com/bar",
		Title:       "Bar",
		PublishedAt: time.Now(),
	}
	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []feed.Post{feedPostFoo, feedPostBar},
	}

	for _, feedPost := range feedBlog.Posts {
		test.AssertEqual(t, feedPost.Content, "")
	}

	pages := map[string]string{
		feedPostFoo.URL: "content about foo",
		feedPostBar.URL: "content about bar",
	}
	pageFetcher := fetchMock.NewPageFetcher(pages)

	feedBlog, err := feed.Hydrate(feedBlog, pageFetcher)
	test.AssertNilError(t, err)

	for _, feedPost := range feedBlog.Posts {
		want, ok := pages[feedPost.URL]
		if !ok {
			t.Errorf("invalid post URL: %s", feedPost.URL)
			continue
		}

		test.AssertEqual(t, feedPost.Content, want)
	}
}
