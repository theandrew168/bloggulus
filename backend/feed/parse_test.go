package feed_test

import (
	"testing"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/theandrew168/bloggulus/backend/feed"
	mockfeed "github.com/theandrew168/bloggulus/backend/feed/mock"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

func TestNormalizePostURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		blogURL string
		postURL string
		want    string
	}{
		{blogURL: "https://example.com", postURL: "/example", want: "https://example.com/example"},
		{blogURL: "https://example.com/", postURL: "/example", want: "https://example.com/example"},
		{blogURL: "http://example.com", postURL: "/example", want: "http://example.com/example"},
		{blogURL: "http://example.com/", postURL: "/example", want: "http://example.com/example"},
		{blogURL: "example.com", postURL: "/example", want: "https://example.com/example"},
		{blogURL: "example.com/", postURL: "/example", want: "https://example.com/example"},
	}
	for _, tt := range tests {
		got := feed.NormalizePostURL(tt.blogURL, tt.postURL)
		test.AssertEqual(t, got, tt.want)
	}
}

func TestDeterminePublishedAt(t *testing.T) {
	t.Parallel()

	feedUpdatedParsed := time.Now().AddDate(0, 0, -1)
	itemPublishedParsed := time.Now().AddDate(0, 0, -2)
	now := time.Now()

	tests := []struct {
		feedUpdatedParsed   *time.Time
		itemPublishedParsed *time.Time
		now                 time.Time
		want                time.Time
	}{
		{
			// Without any feed or item published date, use now.
			feedUpdatedParsed:   nil,
			itemPublishedParsed: nil,
			now:                 now,
			want:                timeutil.Normalize(now),
		},
		{
			// If the feed has an updated date, use it.
			feedUpdatedParsed:   &feedUpdatedParsed,
			itemPublishedParsed: nil,
			now:                 now,
			want:                timeutil.Normalize(feedUpdatedParsed),
		},
		{
			// If the item has a published date, use it.
			feedUpdatedParsed:   nil,
			itemPublishedParsed: &itemPublishedParsed,
			now:                 now,
			want:                timeutil.Normalize(itemPublishedParsed),
		},
		{
			// If item has a published date, use it even if the feed has an updated date.
			feedUpdatedParsed:   &feedUpdatedParsed,
			itemPublishedParsed: &itemPublishedParsed,
			now:                 now,
			want:                timeutil.Normalize(itemPublishedParsed),
		},
	}

	for _, tt := range tests {
		got := feed.DeterminePublishedAt(
			&gofeed.Feed{UpdatedParsed: tt.feedUpdatedParsed},
			&gofeed.Item{PublishedParsed: tt.itemPublishedParsed},
			tt.now,
		)
		test.AssertEqual(t, got, tt.want)
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	feedPostFoo := mockfeed.Post{
		URL:         "https://example.com/foo",
		Title:       "Foo",
		Content:     "content about foo",
		PublishedAt: time.Now(),
	}
	feedPostBar := mockfeed.Post{
		URL:         "https://example.com/bar",
		Title:       "Bar",
		Content:     "content about bar",
		PublishedAt: time.Now(),
	}
	feedBlog := mockfeed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []mockfeed.Post{feedPostFoo, feedPostBar},
	}

	atomFeed, err := mockfeed.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	parsedBlog, err := feed.Parse(test.MustNewURL("https://example.com/atom.xml"), atomFeed)
	test.AssertNilError(t, err)
	test.AssertEqual(t, parsedBlog.Title.Value(), feedBlog.Title)
	test.AssertEqual(t, parsedBlog.SiteURL.Value(), feedBlog.SiteURL)
	test.AssertEqual(t, parsedBlog.FeedURL.Value(), feedBlog.FeedURL)
	test.AssertEqual(t, len(parsedBlog.Posts), len(feedBlog.Posts))

	postsByURL := map[string]mockfeed.Post{
		feedPostFoo.URL: feedPostFoo,
		feedPostBar.URL: feedPostBar,
	}
	for _, parsedPost := range parsedBlog.Posts {
		post, ok := postsByURL[parsedPost.URL.Value()]
		if !ok {
			t.Errorf("invalid post URL: %s", parsedPost.URL.Value())
			continue
		}

		test.AssertEqual(t, parsedPost.URL.Value(), post.URL)
		test.AssertEqual(t, parsedPost.Title.Value(), post.Title)
		test.AssertEqual(t, parsedPost.Content, post.Content)
	}
}

func TestParseMissingURL(t *testing.T) {
	t.Parallel()

	feedPostFoo := mockfeed.Post{
		Title:       "Foo",
		Content:     "content about foo",
		PublishedAt: time.Now(),
	}
	feedBlog := mockfeed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []mockfeed.Post{feedPostFoo},
	}

	atomFeed, err := mockfeed.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	parsedBlog, err := feed.Parse(test.MustNewURL("https://example.com/atom.xml"), atomFeed)
	test.AssertNilError(t, err)

	test.AssertEqual(t, len(parsedBlog.Posts), 0)
}

func TestParseMissingTitle(t *testing.T) {
	t.Parallel()

	feedPostFoo := mockfeed.Post{
		URL:         "https://example.com/foo",
		Content:     "content about foo",
		PublishedAt: time.Now(),
	}
	feedBlog := mockfeed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []mockfeed.Post{feedPostFoo},
	}

	atomFeed, err := mockfeed.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	parsedBlog, err := feed.Parse(test.MustNewURL("https://example.com/atom.xml"), atomFeed)
	test.AssertNilError(t, err)

	test.AssertEqual(t, len(parsedBlog.Posts), 0)
}

func TestParseMissingDomain(t *testing.T) {
	t.Parallel()

	feedPostFoo := mockfeed.Post{
		URL:         "/foo",
		Title:       "Foo",
		Content:     "content about foo",
		PublishedAt: time.Now(),
	}
	feedBlog := mockfeed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []mockfeed.Post{feedPostFoo},
	}

	atomFeed, err := mockfeed.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	parsedBlog, err := feed.Parse(test.MustNewURL("https://example.com/atom.xml"), atomFeed)
	test.AssertNilError(t, err)

	for _, parsedPost := range parsedBlog.Posts {
		test.AssertEqual(t, parsedPost.URL.Value(), feedBlog.SiteURL+feedPostFoo.URL)
	}
}

func TestParseMissingScheme(t *testing.T) {
	t.Parallel()

	feedPostFoo := mockfeed.Post{
		URL:         "example.com/foo",
		Title:       "Foo",
		Content:     "content about foo",
		PublishedAt: time.Now(),
	}
	feedBlog := mockfeed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []mockfeed.Post{feedPostFoo},
	}

	atomFeed, err := mockfeed.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	parsedBlog, err := feed.Parse(test.MustNewURL("https://example.com/atom.xml"), atomFeed)
	test.AssertNilError(t, err)

	for _, parsedPost := range parsedBlog.Posts {
		test.AssertEqual(t, parsedPost.URL.Value(), "https://"+feedPostFoo.URL)
	}
}

func TestParsePublishedAtUTC(t *testing.T) {
	t.Parallel()

	publishedAt, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
	test.AssertNilError(t, err)

	feedPostFoo := mockfeed.Post{
		URL:         "https://example.com/foo",
		Title:       "Foo",
		Content:     "content about foo",
		PublishedAt: publishedAt,
	}
	feedBlog := mockfeed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []mockfeed.Post{feedPostFoo},
	}

	atomFeed, err := mockfeed.GenerateAtomFeed(feedBlog)
	test.AssertNilError(t, err)

	parsedBlog, err := feed.Parse(test.MustNewURL("https://example.com/atom.xml"), atomFeed)
	test.AssertNilError(t, err)

	for _, parsedPost := range parsedBlog.Posts {
		test.AssertEqual(t, parsedPost.PublishedAt, timeutil.Normalize(publishedAt))
	}
}

func BenchmarkParse(b *testing.B) {
	feedPostFoo := mockfeed.Post{
		URL:         "https://example.com/foo",
		Title:       "Foo",
		Content:     "content about foo",
		PublishedAt: time.Now(),
	}
	feedPostBar := mockfeed.Post{
		URL:         "https://example.com/bar",
		Title:       "Bar",
		Content:     "content about bar",
		PublishedAt: time.Now(),
	}
	feedBlog := mockfeed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []mockfeed.Post{feedPostFoo, feedPostBar},
	}

	atomFeed, err := mockfeed.GenerateAtomFeed(feedBlog)
	if err != nil {
		b.Fatalf("got: %v; want: nil", err)
	}

	for b.Loop() {
		feed.Parse(test.MustNewURL("https://example.com/atom.xml"), atomFeed)
	}
}
