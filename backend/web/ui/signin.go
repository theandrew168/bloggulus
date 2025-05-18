package ui

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SignInPageData struct {
	PageLayoutData

	EnableDebugAuth bool
}

func SignInPage(data SignInPageData) g.Node {
	return PageLayout(data.PageLayoutData,
		h.Section(h.Class("signin"),
			h.Article(h.Class("signin__card"),
				h.H2(h.Class("signin__heading"),
					g.Text("Welcome!"),
				),

				h.Hr(),

				g.If(data.EnableDebugAuth,
					h.Form(
						h.Method("POST"),
						h.Action("/debug/signin"),
						h.Input(
							h.Type("hidden"),
							h.Name("csrf_token"),
							h.Value(data.CSRFToken),
						),
						h.Button(h.Class("signin__button"),
							h.Type("submit"),
							h.Img(h.Class("signin__icon"),
								h.Src("/img/bloggulus.png"),
							),
							g.Text("Sign in with Debug"),
						),
					),
				),

				h.A(h.Class("signin__button"),
					h.Href("/github/signin"),
					h.Img(h.Class("signin__icon"),
						h.Src("/img/github.png"),
					),
					g.Text("Sign in with GitHub"),
				),

				h.A(h.Class("signin__button"),
					h.Href("/google/signin"),
					h.Img(h.Class("signin__icon"),
						h.Src("/img/google.png"),
					),
					g.Text("Sign in with Google"),
				),
			),
		),
	)
}
