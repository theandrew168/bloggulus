package ui

import (
	"fmt"

	"github.com/theandrew168/bloggulus/backend/model"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AccountsPageData struct {
	PageLayoutData

	Accounts []*model.Account
}

func AccountsPage(data AccountsPageData) g.Node {
	return PageLayout(data.PageLayoutData,
		h.Section(h.Class("accounts"),
			h.Header(h.Class("accounts-header"),
				h.H1(h.Class("accounts-header__title"),
					g.Text("Accounts"),
				),
			),
			h.Ul(h.Class("accounts-list"),
				h.ID("accounts"),
				g.Map(data.Accounts, func(account *model.Account) g.Node {
					return h.Li(h.Class("accounts-list__item"),
						h.P(
							g.Text(account.Username()),
						),
						h.Form(
							h.Method("POST"),
							h.Action(fmt.Sprintf("/accounts/%s/delete", account.ID())),
							h.Input(
								h.Type("hidden"),
								h.Name("csrf_token"),
								h.Value(data.CSRFToken),
							),
							h.Button(h.Class("button button--outline"),
								h.Type("submit"),
								g.Text("Delete"),
							),
						),
					)
				}),
			),
		),
	)
}
