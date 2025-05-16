package ui

import (
	"github.com/theandrew168/bloggulus/backend/model"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PostPageData struct {
	PageLayoutData

	Post *model.Post
}

func PostPage() g.Node {
	return h.Div(g.Text("TODO: Post"))
}
