package feed

import (
	"log/slog"
	"regexp"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/theandrew168/bloggulus/backend/fetch"
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

		url := item.Link

		// ensure url includes the site's domain
		if url[0] == '/' {
			url = feed.Link + url
		}

		// ensure url includes a scheme (assume https if necessary)
		hasScheme := schemeRegexp.MatchString(url)
		if !hasScheme {
			url = "https://" + url
		}

		// check for a publish date, else default to now
		publishedAt := time.Now()
		if item.PublishedParsed != nil {
			publishedAt = *item.PublishedParsed
		}

		// ensure publishedAt is in UTC
		publishedAt = publishedAt.UTC()

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

func Hydrate(blog Blog, pageFetcher fetch.PageFetcher) (Blog, error) {
	var hydratedPosts []Post
	for _, post := range blog.Posts {
		if post.Content == "" {
			content, err := pageFetcher.FetchPage(post.URL)
			if err != nil {
				slog.Warn("failed to fetch page", "url", post.URL)
			} else {
				post.Content = content
			}
		}
		hydratedPosts = append(hydratedPosts, post)
	}

	blog.Posts = hydratedPosts
	return blog, nil
}
