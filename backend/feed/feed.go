package feed

import (
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

var (
	schemeRegexp = regexp.MustCompile("https?://")
)

type Blog struct {
	FeedURL string
	SiteURL string
	Title   string
	Posts   []Post
}

type Post struct {
	URL         string
	Title       string
	Content     string
	PublishedAt time.Time
}

// Normalize post URLs, ensuring they are full URLs with valid schemes.
func NormalizePostURL(blogURL, postURL string) string {
	url := postURL

	// Ensure post URLs include the site's domain.
	if strings.HasPrefix(postURL, "/") {
		// Omit any duplicate slashes when joining the URLs.
		if strings.HasSuffix(blogURL, "/") {
			url = blogURL + postURL[1:]
		} else {
			url = blogURL + postURL
		}
	}

	// Ensure post URLs include a scheme (assume https if necessary).
	hasScheme := schemeRegexp.MatchString(url)
	if !hasScheme {
		url = "https://" + url
	}

	return url
}

func Parse(feedURL string, feedBody string) (Blog, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseString(feedBody)
	if err != nil {
		return Blog{}, err
	}

	var posts []Post
	for _, item := range feed.Items {
		// skip items without a link or title
		if item.Link == "" || item.Title == "" {
			continue
		}

		url := NormalizePostURL(feed.Link, item.Link)

		// check for a publish date, else default to now
		publishedAt := time.Now()
		if item.PublishedParsed != nil {
			publishedAt = *item.PublishedParsed
		}

		// ensure publishedAt is in UTC
		publishedAt = publishedAt.UTC().Round(time.Microsecond)

		post := Post{
			URL:         url,
			Title:       item.Title,
			Content:     item.Content,
			PublishedAt: publishedAt,
		}
		posts = append(posts, post)
	}

	blog := Blog{
		FeedURL: feedURL,
		SiteURL: feed.Link,
		Title:   feed.Title,
		Posts:   posts,
	}
	return blog, nil
}
