package feed

import (
	"net/url"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/theandrew168/bloggulus/internal/core"
)

// TODO: consume an io.Reader?
func ReadBlog(feedURL string) (core.Blog, error) {
	// early check to ensure the URL is valid
	_, err := url.Parse(feedURL)
	if err != nil {
		return core.Blog{}, err
	}

	// attempt to parse the feed via gofeed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return core.Blog{}, err
	}

	// create a Blog core for the feed
	blog := core.Blog{
		FeedURL: feedURL,
		SiteURL: feed.Link,
		Title:   feed.Title,
	}

	return blog, nil
}

// TODO: consume an io.Reader?
func ReadPosts(blog core.Blog) ([]core.Post, error) {
	// attempt to parse the feed via gofeed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(blog.FeedURL)
	if err != nil {
		return nil, err
	}

	// create a Post core for each entry
	var posts []core.Post
	for _, item := range feed.Items {
		// try Updated then Published to obtain a timestamp
		var updated time.Time
		if item.UpdatedParsed != nil {
			updated = *item.UpdatedParsed
		} else if item.PublishedParsed != nil {
			updated = *item.PublishedParsed
		} else {
			// else default to a month ago
			updated = time.Now().AddDate(0, -1, 0)
		}

		post := core.Post{
			URL:     item.Link,
			Title:   item.Title,
			Updated: updated,
			Blog:    blog,
		}
		posts = append(posts, post)
	}

	return posts, nil
}
