package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed pages.html
var PagesHTML string

type PagesData struct {
	layout.BaseData

	Pages []*model.Page
}

type PagesPage struct {
	tmpl *template.Template
}

func NewPages() *PagesPage {
	sources := []string{
		layout.BaseHTML,
		PagesHTML,
	}

	tmpl := newTemplate("default", sources)
	page := PagesPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *PagesPage) Render(w io.Writer, data PagesData) error {
	return p.tmpl.ExecuteTemplate(w, "default", data)
}

func (p *PagesPage) RenderPages(w io.Writer, data PagesData) error {
	return p.tmpl.ExecuteTemplate(w, "pages", data)
}
