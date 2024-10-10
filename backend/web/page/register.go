package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed register.html
var RegisterHTML string

type RegisterData struct {
	layout.BaseData

	Username string
	Errors   map[string]string
}

type RegisterPage struct {
	tmpl *template.Template
}

func NewRegister() *RegisterPage {
	sources := []string{
		layout.BaseHTML,
		RegisterHTML,
	}

	tmpl := newTemplate("default", sources)
	page := RegisterPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *RegisterPage) Render(w io.Writer, data RegisterData) error {
	return p.tmpl.ExecuteTemplate(w, "default", data)
}
