package fetch

import "errors"

var (
	ErrUnreachablePage = errors.New("fetch: unreachable page")
)

type PageFetcher interface {
	FetchPage(url string) (string, error)
}
