package feed

import (
	"net/url"

	"github.com/mmcdole/gofeed"

	"github.com/theandrew168/bloggulus/internal/core"
)

func ReadBlog(feedURL string) (*core.Blog, error) {
	// early check to ensure the URL is valid
	URL, err := url.Parse(feedURL)
	if err != nil {
		return nil, err
	}

	// use scheme + hostname as site URL
	var siteURL string
	if URL.Scheme != "" {
		siteURL += URL.Scheme + "://"
	}
	siteURL += URL.Hostname()

	// attempt to parse the feed via gofeed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, err
	}

	// create a Blog core for the feed
	blog := core.Blog{
		FeedURL: feedURL,
		SiteURL: siteURL,
		Title:   feed.Title,
	}

	return &blog, nil
}

func ReadPosts(feedURL string) ([]*core.Post, error) {
	// early check to ensure the URL is valid
	_, err := url.Parse(feedURL)
	if err != nil {
		return nil, err
	}

	// attempt to parse the feed via gofeed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, err
	}

	// create a Post core for each entry
	var posts []*core.Post
	for _, item := range feed.Items {
		// try Updated then Published to obtain a timestamp
		updated := item.UpdatedParsed
		if updated == nil {
			updated = item.PublishedParsed
		}
		if updated == nil {
			// TODO: handle dateless posts?
			continue
		}

		post := core.Post{
			URL:     item.Link,
			Title:   item.Title,
			Updated: *updated,
		}
		posts = append(posts, &post)
	}

	return posts, nil
}
