package page

import (
	_ "embed"

	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed blogs.html
var BlogsHTML string

type BlogsData struct {
	layout.BaseData

	Blogs []finder.BlogForAccount
}
