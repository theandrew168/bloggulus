package feed

import "errors"

var (
	ErrUnreachableFeed = errors.New("feed: unreachable feed")
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
