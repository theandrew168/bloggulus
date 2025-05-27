package web

import (
	"io"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/feed"
)

const UserAgent = "Bloggulus/0.5.2 (+https://bloggulus.com)"

// ensure FeedFetcher interface is satisfied
var _ feed.FeedFetcher = (*FeedFetcher)(nil)

type FeedFetcher struct{}

func NewFeedFetcher() *FeedFetcher {
	f := FeedFetcher{}
	return &f
}

func (f *FeedFetcher) FetchFeed(request feed.FetchFeedRequest) (feed.FetchFeedResponse, error) {
	req, err := http.NewRequest("GET", request.URL, nil)
	if err != nil {
		return feed.FetchFeedResponse{}, feed.ErrUnreachableFeed
	}
	req.Header.Set("User-Agent", UserAgent)

	if request.ETag != "" {
		req.Header.Set("If-None-Match", request.ETag)
	}
	if request.LastModified != "" {
		req.Header.Set("If-Modified-Since", request.LastModified)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return feed.FetchFeedResponse{}, feed.ErrUnreachableFeed
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return feed.FetchFeedResponse{}, feed.ErrUnreachableFeed
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return feed.FetchFeedResponse{}, feed.ErrUnreachableFeed
	}

	fetchFeedResponse := feed.FetchFeedResponse{
		Feed:         string(body),
		ETag:         resp.Header.Get("ETag"),
		LastModified: resp.Header.Get("Last-Modified"),
	}

	return fetchFeedResponse, nil
}
