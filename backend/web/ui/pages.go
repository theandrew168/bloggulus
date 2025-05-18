package ui

import (
	"fmt"

	"github.com/theandrew168/bloggulus/backend/model"

	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

type PagesPageData struct {
	PageLayoutData

	Pages []*model.Page
}

func PagesPage(data PagesPageData) g.Node {
	return PageLayout(data.PageLayoutData,
		h.Section(h.Class("pages"),
			h.Header(h.Class("pages-header"),
				h.H1(h.Class("pages-header__title"),
					g.Text("Pages"),
				),
				h.Form(
					h.Method("POST"),
					h.Action("/pages/create"),
					h.Input(
						h.Type("hidden"),
						h.Name("csrf_token"),
						h.Value(data.CSRFToken),
					),
					h.Input(h.Class("input pages-header__input"),
						h.Type("text"),
						h.Name("url"),
						h.Placeholder("Add Page"),
					),
					h.Button(h.Class("button"),
						h.Type("submit"),
						g.Text("Add"),
					),
				),
			),
			PagesList(data),
		),
	)
}

func PagesList(data PagesPageData) g.Node {
	return h.Ul(h.Class("pages-list"),
		h.ID("pages"),
		g.Map(data.Pages, func(page *model.Page) g.Node {
			return h.Li(h.Class("pages-list__item"),
				h.A(h.Class("pages-list__link"),
					h.Href(page.URL()),
					g.Text(page.Title()),
				),
				h.Form(
					h.Method("POST"),
					h.Action(fmt.Sprintf("/pages/%s/unfollow", page.ID())),
					hx.Post(fmt.Sprintf("/pages/%s/unfollow", page.ID())),
					hx.Target("#pages"),
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
	)
}
