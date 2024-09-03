package page

import (
	_ "embed"

	"github.com/theandrew168/bloggulus/backend/finder"
)

//go:embed index.html
var IndexHTML string

type IndexData struct {
	Search       string
	Articles     []finder.Article
	HasMorePages bool
	NextPage     int
}
