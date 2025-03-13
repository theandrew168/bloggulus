package web

import (
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"

	"github.com/theandrew168/bloggulus/backend/fetch"
)

var _ fetch.PageFetcher = (*PageFetcher)(nil)

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

func (f *PageFetcher) FetchPage(request fetch.FetchPageRequest) (fetch.FetchPageResponse, error) {
	req, err := http.NewRequest("GET", request.URL, nil)
	if err != nil {
		return fetch.FetchPageResponse{}, fetch.ErrUnreachablePage
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fetch.FetchPageResponse{}, fetch.ErrUnreachablePage
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fetch.FetchPageResponse{}, fetch.ErrUnreachablePage
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return fetch.FetchPageResponse{}, fetch.ErrUnreachablePage
	}

	response := fetch.FetchPageResponse{
		Content: cleanHTML(string(buf)),
	}
	return response, nil
}
