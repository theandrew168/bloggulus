package feeds

import (
	"log"
	"net/url"

	"github.com/mmcdole/gofeed"

	"github.com/theandrew168/bloggulus/models"
)

func ReadBlog(feedURL string) (*models.Blog, error) {
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

	// create a Blog model for the feed
	blog := models.Blog{
		FeedURL: feedURL,
		SiteURL: siteURL,
		Title:   feed.Title,
	}

	return &blog, nil
}

func ReadPosts(feedURL string) ([]*models.Post, error) {
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

	// create a Post model for each entry
	var posts []*models.Post
	for _, item := range feed.Items {
		// try Updated then Published to obtain a timestamp
		updated := item.UpdatedParsed
		if updated == nil {
			updated = item.PublishedParsed
		}
		if updated == nil {
			log.Printf("skipping dateless post: %s\n", item.Title)
			continue
		}

		post := models.Post{
			URL:     item.Link,
			Title:   item.Title,
			Updated: *updated,
		}
		posts = append(posts, &post)
	}

	return posts, nil
}
