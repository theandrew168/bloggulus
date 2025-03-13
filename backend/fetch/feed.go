package fetch

import "errors"

var (
	ErrUnreachableFeed = errors.New("fetch: unreachable feed")
)

type FetchFeedRequest struct {
	URL          string
	ETag         string
	LastModified string
}

type FetchFeedResponse struct {
	Feed         string
	ETag         string
	LastModified string
}

type FeedFetcher interface {
	FetchFeed(request FetchFeedRequest) (FetchFeedResponse, error)
}
