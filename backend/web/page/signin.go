package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed signin.html
var SignInHTML string

type SignInData struct {
	layout.BaseData

	EnableDebugAuth bool
}

type SignInPage struct {
	tmpl *template.Template
}

func NewSignIn() *SignInPage {
	sources := []string{
		layout.BaseHTML,
		SignInHTML,
	}

	tmpl := newTemplate("default", sources)
	page := SignInPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *SignInPage) Render(w io.Writer, data SignInData) error {
	return p.tmpl.ExecuteTemplate(w, "default", data)
}
