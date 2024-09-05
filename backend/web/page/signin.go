package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed signin.html
var SigninHTML string

type SigninData struct {
	layout.BaseData

	Username string
	Errors   map[string]string
}

type SigninPage struct {
	tmpl *template.Template
}

func NewSignin() *SigninPage {
	sources := []string{
		layout.BaseHTML,
		SigninHTML,
	}

	tmpl := newTemplate("page", sources)
	page := SigninPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *SigninPage) Render(w io.Writer, data SigninData) error {
	return p.tmpl.ExecuteTemplate(w, "page", data)
}
