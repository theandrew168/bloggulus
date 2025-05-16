package ui

import (
	"fmt"

	"github.com/theandrew168/bloggulus/backend/finder"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type IndexPageData struct {
	PageLayoutData

	Search       string
	Articles     []finder.Article
	HasMorePages bool
	NextPage     int
}

func IndexPage(data IndexPageData) g.Node {
	return PageLayout(data.PageLayoutData,
		// Header (search)
		h.Header(
			h.Class("articles-header"),
			g.If(data.Search != "", h.H1(
				h.Class("articles-header__title"),
				g.Text("Relevant Articles"),
			)),
			g.If(data.Search == "", h.H1(
				h.Class("articles-header__title"),
				g.Text("Recent Articles"),
			)),
			// TODO: Replace this with <search> once gomponents supports it.
			h.Div(
				h.Form(
					h.Method("GET"),
					h.Action("/"),
					h.Input(
						h.Class("input"),
						h.Type("text"),
						h.Name("q"),
						h.Value(data.Search),
						h.Placeholder("Search"),
					),
				),
			),
		),

		// Articles
		h.Section(
			h.Class("articles"),
			g.Map(data.Articles, func(article finder.Article) g.Node {
				return h.Article(
					h.Class("article"),
					h.Header(
						h.Class("article__header"),
						h.Span(
							h.Class("article__date"),
							g.Text(article.PublishedAt.Format("Jan 2, 2006")),
						),
						h.Ul(
							h.Class("article__tags"),
							g.Map(article.Tags, func(tag string) g.Node {
								return h.Li(
									h.A(
										h.Class("article__tag"),
										h.Href(fmt.Sprintf("/?q=%s", tag)),
										g.Text(tag),
									),
								)
							}),
						),
					),
					h.P(
						h.A(
							h.Class("article__title"),
							h.Href(article.URL),
							g.Text(article.Title),
						),
					),
					h.P(
						h.A(
							h.Class("article__blog-title"),
							h.Href(article.BlogURL),
							g.Text(article.BlogTitle),
						),
					),
				)
			}),
			g.If(len(data.Articles) == 0,
				g.Group{
					g.If(data.Search != "",
						h.Article(
							h.Class("articles-cta"),
							h.P(
								g.Text("No relevant articles! Try searching for something else."),
							),
						),
					),
					g.If(data.Search == "",
						h.Article(
							h.Class("articles-cta"),
							h.P(
								g.Text("No posts found! Get started by following your favorite blogs."),
							),
							h.P(
								h.A(
									h.Class("button"),
									h.Href("/blogs"),
									g.Text("Follow Blogs"),
								),
							),
						),
					),
				},
			),
		),

		// Footer (pagination)
		g.If(data.HasMorePages,
			h.Footer(
				h.Class("articles-footer"),
				g.If(data.Search == "", h.A(
					h.Class("button button--outline"),
					h.Href(fmt.Sprintf("/?p=%d&q=%s", data.NextPage, data.Search)),
					g.Text("See More"),
				)),
				g.If(data.Search != "", h.A(
					h.Class("button button--outline"),
					h.Href(fmt.Sprintf("/?p=%d", data.NextPage)),
					g.Text("See More"),
				)),
			),
		),
	)
}
