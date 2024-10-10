package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed error.html
var ErrorHTML string

type ErrorData struct {
	layout.BaseData

	StatusCode int
	StatusText string
	Message    string
}

type ErrorPage struct {
	tmpl *template.Template
}

func NewError() *ErrorPage {
	sources := []string{
		layout.BaseHTML,
		ErrorHTML,
	}

	tmpl := newTemplate("default", sources)
	page := ErrorPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *ErrorPage) Render(w io.Writer, data ErrorData) error {
	return p.tmpl.ExecuteTemplate(w, "default", data)
}
