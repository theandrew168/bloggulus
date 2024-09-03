package page

import (
	_ "embed"

	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed index.html
var IndexHTML string

type IndexData struct {
	layout.BaseData

	Search       string
	Articles     []finder.Article
	HasMorePages bool
	NextPage     int
}
