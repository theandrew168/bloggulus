package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/web/layout"
	"golang.org/x/oauth2"
)

//go:embed login.html
var LoginHTML string

type LoginData struct {
	layout.BaseData

	GithubConf *oauth2.Config
	Errors     map[string]string

	EnableDebugAuth bool
}

type LoginPage struct {
	tmpl *template.Template
}

func NewLogin() *LoginPage {
	sources := []string{
		layout.BaseHTML,
		LoginHTML,
	}

	tmpl := newTemplate("default", sources)
	page := LoginPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *LoginPage) Render(w io.Writer, data LoginData) error {
	return p.tmpl.ExecuteTemplate(w, "default", data)
}
