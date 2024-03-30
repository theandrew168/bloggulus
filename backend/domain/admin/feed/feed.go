package feed

import (
	"fmt"
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
	Contents    string
	PublishedAt time.Time
}

// TODO: Better, richer error handling that includes data about
// errors from individual posts (invalid URLs, for example).
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

func Hydrate(blog Blog, pageFetcher PageFetcher) (Blog, error) {
	var hydratedPosts []Post
	for _, post := range blog.Posts {
		if post.Contents == "" {
			content, err := pageFetcher.FetchPage(post.URL)
			if err != nil {
				return Blog{}, err
			}

			post.Contents = content
		}
		hydratedPosts = append(hydratedPosts, post)
	}

	blog.Posts = hydratedPosts
	return blog, nil
}
