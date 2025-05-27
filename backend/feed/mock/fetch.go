package mock

import (
	"github.com/theandrew168/bloggulus/backend/feed"
)

// ensure FeedFetcher interface is satisfied
var _ feed.FeedFetcher = (*FeedFetcher)(nil)

type FeedFetcher struct {
	feeds map[string]feed.FetchFeedResponse
}

func NewFeedFetcher(feeds map[string]feed.FetchFeedResponse) *FeedFetcher {
	f := FeedFetcher{feeds: feeds}
	return &f
}

func (f *FeedFetcher) FetchFeed(request feed.FetchFeedRequest) (feed.FetchFeedResponse, error) {
	feedForURL, ok := f.feeds[request.URL]
	if !ok {
		return feed.FetchFeedResponse{}, feed.ErrUnreachableFeed
	}

	return feedForURL, nil
}
