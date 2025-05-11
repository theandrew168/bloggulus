package ui

import (
	"github.com/theandrew168/bloggulus/backend/model"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type LayoutData struct {
	Account   *model.Account
	CSRFToken string
	Toast     string

	GoatCounterCode     string
	PlausibleDataDomain string
}

func Layout(data LayoutData) Node {
	return HTML5(HTML5Props{
		Language:    "en",
		Title:       "Bloggulus - A website for avid blog readers",
		Description: "Bloggulus - A website for avid blog readers",
		Head:        []Node{},
		Body: []Node{
			Div(Text("wow")),
		},
	})
}
