package feed_test

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/theandrew168/bloggulus/backend/domain/admin/feed"
	"github.com/theandrew168/bloggulus/backend/domain/admin/fetch"
	"github.com/theandrew168/bloggulus/backend/test"
)

// Convert a feed.Blog into an XML (Atom) document.
func generateAtomFeed(t *testing.T, blog feed.Blog) string {
	t.Helper()

	type xmlLink struct {
		HREF string `xml:"href,attr"`
		Rel  string `xml:"rel,attr,omitempty"`
	}

	type xmlPost struct {
		URL         xmlLink   `xml:"link"`
		Title       string    `xml:"title"`
		Contents    string    `xml:"content"`
		PublishedAt time.Time `xml:"published"`
	}

	type xmlBlog struct {
		XMLName xml.Name  `xml:"feed"`
		Links   []xmlLink `xml:"link"`
		Title   string    `xml:"title"`
		Posts   []xmlPost `xml:"entry"`
	}

	var posts []xmlPost
	for _, post := range blog.Posts {
		posts = append(posts, xmlPost{
			URL:         xmlLink{HREF: post.URL},
			Title:       post.Title,
			Contents:    post.Contents,
			PublishedAt: post.PublishedAt,
		})
	}

	b := xmlBlog{
		Links: []xmlLink{
			{HREF: blog.FeedURL, Rel: "self"},
			{HREF: blog.SiteURL, Rel: "alternate"},
		},
		Title: blog.Title,
		Posts: posts,
	}

	out, err := xml.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	return xml.Header + string(out)
}

func TestParse(t *testing.T) {
	t.Parallel()

	feedPostFoo := feed.Post{
		URL:         "https://example.com/foo",
		Title:       "Foo",
		Contents:    "content about foo",
		PublishedAt: time.Now(),
	}
	feedPostBar := feed.Post{
		URL:         "https://example.com/bar",
		Title:       "Bar",
		Contents:    "content about bar",
		PublishedAt: time.Now(),
	}
	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []feed.Post{feedPostFoo, feedPostBar},
	}

	atomFeed := generateAtomFeed(t, feedBlog)

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
		test.AssertEqual(t, parsedPost.Contents, post.Contents)
	}
}

func TestParseMissingDomain(t *testing.T) {
	t.Parallel()

	feedPostFoo := feed.Post{
		URL:         "/foo",
		Title:       "Foo",
		Contents:    "content about foo",
		PublishedAt: time.Now(),
	}
	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []feed.Post{feedPostFoo},
	}

	atomFeed := generateAtomFeed(t, feedBlog)

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
		Contents:    "content about foo",
		PublishedAt: time.Now(),
	}
	feedBlog := feed.Blog{
		Title:   "FooBar",
		SiteURL: "https://example.com",
		FeedURL: "https://example.com/atom.xml",
		Posts:   []feed.Post{feedPostFoo},
	}

	atomFeed := generateAtomFeed(t, feedBlog)

	parsedBlog, err := feed.Parse("https://example.com/atom.xml", atomFeed)
	test.AssertNilError(t, err)

	for _, parsedPost := range parsedBlog.Posts {
		test.AssertEqual(t, parsedPost.URL, "https://"+feedPostFoo.URL)
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
		test.AssertEqual(t, feedPost.Contents, "")
	}

	pages := map[string]string{
		feedPostFoo.URL: "content about foo",
		feedPostBar.URL: "content about bar",
	}
	pageFetcher := fetch.NewMockPageFetcher(pages)

	feedBlog, err := feed.Hydrate(feedBlog, pageFetcher)
	test.AssertNilError(t, err)

	for _, feedPost := range feedBlog.Posts {
		want, ok := pages[feedPost.URL]
		if !ok {
			t.Errorf("invalid post URL: %s", feedPost.URL)
			continue
		}

		test.AssertEqual(t, feedPost.Contents, want)
	}
}
