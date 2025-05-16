package ui

import (
	"github.com/theandrew168/bloggulus/backend/finder"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type BlogsPageData struct {
	PageLayoutData

	Blogs []finder.BlogForAccount
}

func BlogsPage() g.Node {
	return h.Div(g.Text("TODO: Blogs"))
}
