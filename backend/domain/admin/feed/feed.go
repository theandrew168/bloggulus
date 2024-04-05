package feed

import (
	"log/slog"
	"regexp"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/theandrew168/bloggulus/backend/domain/admin/fetch"
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
	Contents    string
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
		// ensure link is valid
		link := item.Link

		// ensure link includes the site's domain
		if link[0] == '/' {
			link = feed.Link + link
		}

		// ensure link includes a scheme (assume https if necessary)
		hasScheme := schemeRegexp.MatchString(link)
		if !hasScheme {
			link = "https://" + link
		}

		// check for a publish date, else default to now
		publishedAt := time.Now()
		if item.PublishedParsed != nil {
			publishedAt = *item.PublishedParsed
		}

		post := Post{
			URL:         link,
			Title:       item.Title,
			Contents:    item.Content,
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
		if post.Contents == "" {
			content, err := pageFetcher.FetchPage(post.URL)
			if err != nil {
				slog.Warn("failed to fetch page", "url", post.URL)
				continue
			}

			post.Contents = content
		}
		hydratedPosts = append(hydratedPosts, post)
	}

	blog.Posts = hydratedPosts
	return blog, nil
}
