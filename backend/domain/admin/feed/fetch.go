package feed

import "time"

type PageFetcher interface {
	fetchPage(url string) (string, error)
}

type FetchFeedResponse struct {
	Feed         string
	ETag         string
	LastModified string
}

type FeedFetcher interface {
	fetchFeed(url, etag, lastModified string) (FetchFeedResponse, error)
}

type Blog struct {
	FeedURL string
	SiteURL string
	Title   string
	Posts   []Post
}

type Post struct {
	URL       string
	Title     string
	UpdatedAt time.Time
	Body      string
}

// TODO: implement this
func ParseFeed(feed string) Blog {
	blog := Blog{}
	return blog
}
