package mock

import (
	"github.com/theandrew168/bloggulus/backend/fetch"
)

// ensure PageFetcher interface is satisfied
var _ fetch.PageFetcher = (*PageFetcher)(nil)

type PageFetcher struct {
	pages map[string]string
}

func NewPageFetcher(pages map[string]string) *PageFetcher {
	f := PageFetcher{pages: pages}
	return &f
}

func (f *PageFetcher) FetchPage(url string) (string, error) {
	page, ok := f.pages[url]
	if !ok {
		return "", fetch.ErrUnreachablePage
	}

	return page, nil
}
