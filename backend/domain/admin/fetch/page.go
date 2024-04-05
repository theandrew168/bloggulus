package fetch

import "errors"

type PageFetcher interface {
	FetchPage(url string) (string, error)
}

// ensure PageFetcher interface is satisfied
var _ PageFetcher = (*MockPageFetcher)(nil)

type MockPageFetcher struct {
	pages map[string]string
}

func NewMockPageFetcher(pages map[string]string) *MockPageFetcher {
	f := MockPageFetcher{pages: pages}
	return &f
}

func (f *MockPageFetcher) FetchPage(url string) (string, error) {
	page, ok := f.pages[url]
	if !ok {
		return "", errors.New("page not found")
	}

	return page, nil
}
