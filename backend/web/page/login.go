package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed login.html
var LoginHTML string

type LoginData struct {
	layout.BaseData

	NextPath string
	Username string
	Errors   map[string]string
}

type LoginPage struct {
	tmpl *template.Template
}

func NewLogin() *LoginPage {
	sources := []string{
		layout.BaseHTML,
		LoginHTML,
	}

	tmpl := newTemplate("page", sources)
	page := LoginPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *LoginPage) Render(w io.Writer, data LoginData) error {
	return p.tmpl.ExecuteTemplate(w, "page", data)
}
