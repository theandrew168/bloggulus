package page

import (
	_ "embed"
	"html/template"
	"io"

	webquery "github.com/theandrew168/bloggulus/backend/query/web"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed index.html
var IndexHTML string

type IndexData struct {
	layout.BaseData

	Search       string
	Articles     []webquery.Article
	HasMorePages bool
	NextPage     int
}

type IndexPage struct {
	tmpl *template.Template
}

func NewIndex() *IndexPage {
	sources := []string{
		layout.BaseHTML,
		IndexHTML,
	}

	tmpl := newTemplate("default", sources)
	page := IndexPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *IndexPage) Render(w io.Writer, data IndexData) error {
	data.Stylesheets = []string{
		"/css/article.css",
	}
	return p.tmpl.ExecuteTemplate(w, "default", data)
}
