package ui

import (
	"github.com/theandrew168/bloggulus/backend/model"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PagesPageData struct {
	PageLayoutData

	Pages []*model.Page
}

func PagesPage() g.Node {
	return h.Div(g.Text("TODO: Pages"))
}
