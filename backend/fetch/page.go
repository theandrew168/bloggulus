package fetch

import "errors"

var (
	ErrUnreachablePage = errors.New("fetch: unreachable page")
)

type FetchPageRequest struct {
	URL string
}

type FetchPageResponse struct {
	Content string
}

type PageFetcher interface {
	FetchPage(request FetchPageRequest) (FetchPageResponse, error)
}
