package ui

import (
	"golang.org/x/oauth2"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SignInPageData struct {
	PageLayoutData

	GithubConf *oauth2.Config
	Errors     map[string]string

	EnableDebugAuth bool
}

func SignInPage() g.Node {
	return h.Div(g.Text("TODO: Sign In"))
}
