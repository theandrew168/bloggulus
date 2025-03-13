package mock

import (
	"github.com/theandrew168/bloggulus/backend/fetch"
)

// ensure FeedFetcher interface is satisfied
var _ fetch.FeedFetcher = (*FeedFetcher)(nil)

type FeedFetcher struct {
	feeds map[string]fetch.FetchFeedResponse
}

func NewFeedFetcher(feeds map[string]fetch.FetchFeedResponse) *FeedFetcher {
	f := FeedFetcher{feeds: feeds}
	return &f
}

func (f *FeedFetcher) FetchFeed(request fetch.FetchFeedRequest) (fetch.FetchFeedResponse, error) {
	feed, ok := f.feeds[request.URL]
	if !ok {
		return fetch.FetchFeedResponse{}, fetch.ErrUnreachableFeed
	}

	return feed, nil
}
