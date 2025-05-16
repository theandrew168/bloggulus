package ui

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ErrorPageData struct {
	PageLayoutData

	StatusCode int
	StatusText string
	Message    string
}

func ErrorPage() g.Node {
	return h.Div(g.Text("TODO: Error"))
}
