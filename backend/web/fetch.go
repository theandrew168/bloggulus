package web

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"

	"github.com/theandrew168/bloggulus/backend/domain/admin/feed"
)

// ensure PageFetcher interface is satisfied
var _ feed.PageFetcher = (*PageFetcher)(nil)

// ensure FeedFetcher interface is satisfied
var _ feed.FeedFetcher = (*FeedFetcher)(nil)

// I know...
var (
	codePattern   = regexp.MustCompile(`(?s)<code>.*?</code>`)
	footerPattern = regexp.MustCompile(`(?s)<footer>.*?</footer>`)
	headerPattern = regexp.MustCompile(`(?s)<header>.*?</header>`)
	navPattern    = regexp.MustCompile(`(?s)<nav>.*?</nav>`)
	prePattern    = regexp.MustCompile(`(?s)<pre>.*?</pre>`)
)

// please PR a better way :(
func cleanHTML(s string) string {
	s = codePattern.ReplaceAllString(s, "")
	s = footerPattern.ReplaceAllString(s, "")
	s = headerPattern.ReplaceAllString(s, "")
	s = navPattern.ReplaceAllString(s, "")
	s = prePattern.ReplaceAllString(s, "")

	s = bluemonday.StrictPolicy().Sanitize(s)
	s = html.UnescapeString(s)
	s = strings.ToValidUTF8(s, "")

	return s
}

type PageFetcher struct{}

func NewPageFetcher() *PageFetcher {
	f := PageFetcher{}
	return &f
}

func (f *PageFetcher) FetchPage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("%v: %v", url, err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%v: %v", url, err)
	}

	return cleanHTML(string(buf)), nil
}

type FeedFetcher struct{}

func NewFeedFetcher() *FeedFetcher {
	f := FeedFetcher{}
	return &f
}

func (f *FeedFetcher) FetchFeed(url, etag, lastModified string) (feed.FetchFeedResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return feed.FetchFeedResponse{}, fmt.Errorf("%v: %v", url, err)
	}

	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	if lastModified != "" {
		req.Header.Set("If-Modified-Since", lastModified)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return feed.FetchFeedResponse{}, fmt.Errorf("%v: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return feed.FetchFeedResponse{}, nil
	}

	fetchFeedResponse := feed.FetchFeedResponse{
		Feed:         string(body),
		ETag:         resp.Header.Get("ETag"),
		LastModified: resp.Header.Get("LastModified"),
	}

	return fetchFeedResponse, nil
}
