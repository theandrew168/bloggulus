package ui

import (
	"fmt"

	"github.com/theandrew168/bloggulus/backend/model"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type PageLayoutData struct {
	Account   *model.Account
	CSRFToken string
	Toast     string

	GoatCounterCode     string
	PlausibleDataDomain string
}

func PageLayout(data PageLayoutData, children ...g.Node) g.Node {
	return c.HTML5(c.HTML5Props{
		Language:    "en",
		Title:       "Bloggulus - A website for avid blog readers",
		Description: "Bloggulus - A website for avid blog readers",
		Head: []g.Node{
			h.Link(h.Rel("stylesheet"), h.Href("/css/reset.css")),
			h.Link(h.Rel("stylesheet"), h.Href("/css/fonts.css")),
			h.Link(h.Rel("stylesheet"), h.Href("/css/style.css")),
			h.Script(h.Src("/js/htmx.min.js"), h.Defer()),
			g.If(data.GoatCounterCode != "", h.Script(
				h.Src("//gc.zgo.at/count.js"),
				g.Attr("data-goatcounter", fmt.Sprintf("https://%s.goatcounter.com/count", data.GoatCounterCode)),
				h.Async(),
			)),
		},
		Body: []g.Node{h.Class("sans-serif"),
			// Header (nav bar)
			h.Header(h.Class("header"),
				h.Nav(
					h.Ul(h.Class("header__links"),
						h.Li(h.Class("header__link--first"),
							h.A(h.Class("header__link header__link--home"),
								h.Href("/"),
								g.Text("Bloggulus"),
							),
						),
						g.If(data.Account != nil,
							g.Group{
								h.Li(
									h.A(h.Class("header__link"),
										h.Href("/blogs"),
										g.Text("Blogs"),
									),
								),
								h.Li(
									h.Form(
										h.Method("POST"),
										h.Action("/signout"),
										h.Input(
											h.Type("hidden"),
											h.Name("csrf_token"),
											h.Value(data.CSRFToken),
										),
										h.Button(h.Class("header__link"),
											h.Type("submit"),
											g.Text("Sign Out"),
										),
									),
								),
							},
						),
						g.If(data.Account == nil,
							h.Li(
								h.A(h.Class("header__link"),
									h.Href("/signin"),
									g.Text("Sign In"),
								),
							),
						),
					),
				),
			),
			// Toast message
			g.If(data.Toast != "",
				h.Div(h.Class("toast"),
					h.Div(
						g.Text(data.Toast),
					),
				),
			),
			// Main content
			h.Main(
				g.Group(children),
			),
			// Footer
			h.Footer(h.Class("footer"),
				h.Nav(
					h.Ul(h.Class("footer__links"),
						h.Li(
							h.A(h.Class("footer__link footer__link--big"),
								h.Href("/"),
								g.Text("Bloggulus"),
							),
						),
						h.Li(
							h.A(h.Class("footer__link footer__link--small"),
								h.Href("/docs/privacy.html"),
								g.Text("Privacy Policy"),
							),
						),
						h.Li(
							h.A(h.Class("footer__link"),
								h.Href("https://shallowbrooksoftware.com"),
								g.Text("Shallow Brook Software"),
							),
						),
					),
				),
			),
		},
	})
}
