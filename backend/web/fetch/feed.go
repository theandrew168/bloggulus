package fetch

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/domain/admin/fetch"
)

var _ fetch.FeedFetcher = (*FeedFetcher)(nil)

type FeedFetcher struct{}

func NewFeedFetcher() *FeedFetcher {
	f := FeedFetcher{}
	return &f
}

func (f *FeedFetcher) FetchFeed(url, etag, lastModified string) (fetch.FetchFeedResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fetch.FetchFeedResponse{}, fmt.Errorf("%v: %v", url, err)
	}

	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	if lastModified != "" {
		req.Header.Set("If-Modified-Since", lastModified)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fetch.FetchFeedResponse{}, fmt.Errorf("%v: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		slog.Info("no new data for feed", "url", url)
		return fetch.FetchFeedResponse{}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fetch.FetchFeedResponse{}, nil
	}

	fetchFeedResponse := fetch.FetchFeedResponse{
		Feed:         string(body),
		ETag:         resp.Header.Get("ETag"),
		LastModified: resp.Header.Get("LastModified"),
	}

	return fetchFeedResponse, nil
}
