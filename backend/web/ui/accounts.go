package ui

import (
	"github.com/theandrew168/bloggulus/backend/model"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AccountsPageData struct {
	PageLayoutData

	Accounts []*model.Account
}

func AccountsPage() g.Node {
	return h.Div(g.Text("TODO: Accounts"))
}
