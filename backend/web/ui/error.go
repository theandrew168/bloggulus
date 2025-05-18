package ui

import (
	"strconv"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ErrorPageData struct {
	PageLayoutData

	StatusCode int
	StatusText string
	Message    string
}

func ErrorPage(data ErrorPageData) g.Node {
	return PageLayout(data.PageLayoutData,
		h.Section(h.Class("error"),
			h.Article(
				h.P(h.Class("error__code"),
					g.Text(strconv.Itoa(data.StatusCode)),
				),
				h.H1(h.Class("error__status"),
					g.Text(data.StatusText),
				),
				h.H2(h.Class("error__message"),
					g.Text(data.Message),
				),
				h.A(h.Class("error__link"),
					h.Href("/"),
					g.Text("Go back home"),
				),
			),
		),
	)
}
