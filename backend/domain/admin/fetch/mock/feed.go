package mock

import (
	"errors"

	"github.com/theandrew168/bloggulus/backend/domain/admin/fetch"
)

// ensure FeedFetcher interface is satisfied
var _ fetch.FeedFetcher = (*FeedFetcher)(nil)

type FeedFetcher struct {
	feeds map[string]string
}

func NewFeedFetcher(feeds map[string]string) *FeedFetcher {
	f := FeedFetcher{feeds: feeds}
	return &f
}

func (f *FeedFetcher) FetchFeed(url, etag, lastModified string) (fetch.FetchFeedResponse, error) {
	feed, ok := f.feeds[url]
	if !ok {
		return fetch.FetchFeedResponse{}, errors.New("feed not found")
	}

	resp := fetch.FetchFeedResponse{
		Feed: feed,
	}
	return resp, nil
}
