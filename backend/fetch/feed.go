package fetch

import "errors"

var (
	ErrUnreachableFeed = errors.New("fetch: unreachable feed")
)

type FetchFeedResponse struct {
	Feed         string
	ETag         string
	LastModified string
}

type FeedFetcher interface {
	FetchFeed(url, etag, lastModified string) (FetchFeedResponse, error)
}
