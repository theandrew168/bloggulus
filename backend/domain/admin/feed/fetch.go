package feed

type PageFetcher interface {
	FetchPage(url string) (string, error)
}

type FetchFeedResponse struct {
	Feed         string
	ETag         string
	LastModified string
}

type FeedFetcher interface {
	FetchFeed(url, etag, lastModified string) (FetchFeedResponse, error)
}
