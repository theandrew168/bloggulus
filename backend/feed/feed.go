package feed

import (
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/theandrew168/bloggulus/backend/timeutil"
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

func DeterminePublishedAt(feed *gofeed.Feed, item *gofeed.Item, now time.Time) time.Time {
	// If the item has a published date, use it since it is the most accurate / specific.
	if item.PublishedParsed != nil {
		return timeutil.Normalize(*item.PublishedParsed)
	}

	// Otherwise, if the feed itself has an updated date, use it instead.
	if feed.UpdatedParsed != nil {
		return timeutil.Normalize(*feed.UpdatedParsed)
	}

	// If all else fails, use the current time.
	return timeutil.Normalize(now)
}

func Parse(feedURL string, feedBody string) (Blog, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseString(feedBody)
	if err != nil {
		return Blog{}, err
	}

	var posts []Post
	for _, item := range feed.Items {
		// Skip items without a link or title.
		if item.Link == "" || item.Title == "" {
			continue
		}

		url := NormalizePostURL(feed.Link, item.Link)
		publishedAt := DeterminePublishedAt(feed, item, time.Now())

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
