package feed

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/microcosm-cc/bluemonday"
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
	blog := core.NewBlog(feedURL, feed.Link, feed.Title)
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

		post := core.NewPost(item.Link, item.Title, updated, blog)
		posts = append(posts, post)
	}

	return posts, nil
}

func ReadPostBody(post core.Post) (string, error) {
	resp, err := http.Get(post.URL)
	if err != nil {
		return "", fmt.Errorf("%v: %v", post.URL, err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%v: %v", post.URL, err)
	}

	p := bluemonday.StripTagsPolicy()
	body := string(buf)
	body = p.Sanitize(body)
	body = html.UnescapeString(body)

	return body, nil
}
