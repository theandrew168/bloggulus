package mock

import (
	"errors"

	"github.com/theandrew168/bloggulus/backend/domain/admin/fetch"
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
		return "", errors.New("page not found")
	}

	return page, nil
}
