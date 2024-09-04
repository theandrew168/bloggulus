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

func NewIndex() (*IndexPage, error) {
	// Create the template.
	tmpl := template.New("page")

	// List all required sources.
	sources := []string{
		layout.BaseHTML,
		IndexHTML,
	}

	// Parse each source required to render this page.
	for _, source := range sources {
		_, err := tmpl.Parse(source)
		if err != nil {
			return nil, err
		}
	}

	page := IndexPage{
		tmpl: tmpl,
	}
	return &page, nil
}

func (p *IndexPage) Render(w io.Writer, data IndexData) error {
	return p.tmpl.ExecuteTemplate(w, "page", data)
}
