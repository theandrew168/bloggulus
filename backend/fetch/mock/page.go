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

func (f *PageFetcher) FetchPage(request fetch.FetchPageRequest) (fetch.FetchPageResponse, error) {
	page, ok := f.pages[request.URL]
	if !ok {
		return fetch.FetchPageResponse{}, fetch.ErrUnreachablePage
	}

	response := fetch.FetchPageResponse{
		Content: page,
	}
	return response, nil
}
