package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed index.html
var IndexHTML string

type IndexData struct {
	layout.BaseData

	Search       string
	Articles     []finder.Article
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

	tmpl := newTemplate("page", sources)
	page := IndexPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *IndexPage) Render(w io.Writer, data IndexData) error {
	return p.tmpl.ExecuteTemplate(w, "page", data)
}
