package feed

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"time"

	"github.com/mmcdole/gofeed"
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
	PublishedAt time.Time
	Body        string
}

// TODO: Better, richer error handling that includes data about
// errors from individual posts (invalid URLs, for example).
func Parse(feedURL string, feedBody io.Reader) (Blog, error) {
	fp := gofeed.NewParser()
	feed, err := fp.Parse(feedBody)
	if err != nil {
		return Blog{}, err
	}

	var posts []Post
	for _, item := range feed.Items {
		// ensure link is valid
		link := item.Link
		u, err := url.Parse(link)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// ensure link includes hostname
		if u.Hostname() == "" {
			link = feed.Link + link
		}

		// ensure link includes scheme
		matched, err := regexp.MatchString("^https?://", link)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// assume https if no scheme is present
		if !matched {
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
			PublishedAt: publishedAt,
			Body:        item.Content,
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

func Hydrate(blog Blog, pageFetcher PageFetcher) error {
	for _, post := range blog.Posts {
		if post.Body == "" {
			body, err := pageFetcher.FetchPage(post.URL)
			if err != nil {
				return err
			}

			post.Body = body
		}
	}

	return nil
}
