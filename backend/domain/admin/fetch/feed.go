package fetch

import "errors"

type FetchFeedResponse struct {
	Feed         string
	ETag         string
	LastModified string
}

type FeedFetcher interface {
	FetchFeed(url, etag, lastModified string) (FetchFeedResponse, error)
}

// ensure FeedFetcher interface is satisfied
var _ FeedFetcher = (*MockFeedFetcher)(nil)

type MockFeedFetcher struct {
	feeds map[string]string
}

func NewMMockFeedFetcher(feeds map[string]string) *MockFeedFetcher {
	f := MockFeedFetcher{feeds: feeds}
	return &f
}

func (f *MockFeedFetcher) FetchFeed(url, etag, lastModified string) (FetchFeedResponse, error) {
	feed, ok := f.feeds[url]
	if !ok {
		return FetchFeedResponse{}, errors.New("feed not found")
	}

	resp := FetchFeedResponse{
		Feed: feed,
	}
	return resp, nil
}
