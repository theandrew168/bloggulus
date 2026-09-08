package feed

import (
	"errors"

	"github.com/theandrew168/bloggulus/backend/value"
)

var (
	ErrUnreachableFeed = errors.New("feed: unreachable feed")
)

type FetchFeedRequest struct {
	URL          value.URL
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
