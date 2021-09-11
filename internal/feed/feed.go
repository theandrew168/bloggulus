package feed

import (
	"net/url"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/theandrew168/bloggulus/internal/core"
)

func ReadBlog(feedURL string) (core.Blog, error) {
	// TODO: check if feed contains its link else default to hostname

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
		SiteURL: feed.Link,  // TODO: check this
		Title:   feed.Title,
	}

	return blog, nil
}

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
		// check if author is present in the feed else default to title
		var author string
		if feed.Author != nil {
			author = feed.Author.Name
		} else {
			author = feed.Title
		}

		// try Updated then Published to obtain a timestamp
		var updated time.Time
		if item.UpdatedParsed != nil {
			updated = *item.UpdatedParsed
		} else if item.PublishedParsed != nil {
			updated = *item.PublishedParsed
		} else {
			// else default to a week ago
			updated = time.Now().AddDate(0, 0, -7)
		}

		// don't worry about Post.Body here (sync_blogs task will check later)
		post := core.Post{
			URL:     item.Link,
			Title:   item.Title,
			Author:  author,
			Updated: updated,
			Blog:    blog,
		}
		posts = append(posts, post)
	}

	return posts, nil
}
