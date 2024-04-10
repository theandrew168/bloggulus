package fetch

type PageFetcher interface {
	FetchPage(url string) (string, error)
}
