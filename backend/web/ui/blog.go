package ui

import (
	"github.com/theandrew168/bloggulus/backend/model"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type BlogPageData struct {
	PageLayoutData

	Blog  *model.Blog
	Posts []*model.Post
}

func BlogPage() g.Node {
	return h.Div(g.Text("TODO: Blog"))
}
