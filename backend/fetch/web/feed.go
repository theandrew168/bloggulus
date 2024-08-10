package web

import (
	"io"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/fetch"
)

// ensure FeedFetcher interface is satisfied
var _ fetch.FeedFetcher = (*FeedFetcher)(nil)

type FeedFetcher struct{}

func NewFeedFetcher() *FeedFetcher {
	f := FeedFetcher{}
	return &f
}

func (f *FeedFetcher) FetchFeed(url, etag, lastModified string) (fetch.FetchFeedResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fetch.FetchFeedResponse{}, fetch.ErrUnreachableFeed
	}
	req.Header.Set("User-Agent", UserAgent)

	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	if lastModified != "" {
		req.Header.Set("If-Modified-Since", lastModified)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fetch.FetchFeedResponse{}, fetch.ErrUnreachableFeed
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fetch.FetchFeedResponse{}, fetch.ErrUnreachableFeed
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fetch.FetchFeedResponse{}, fetch.ErrUnreachablePage
	}

	fetchFeedResponse := fetch.FetchFeedResponse{
		Feed:         string(body),
		ETag:         resp.Header.Get("ETag"),
		LastModified: resp.Header.Get("Last-Modified"),
	}

	return fetchFeedResponse, nil
}
