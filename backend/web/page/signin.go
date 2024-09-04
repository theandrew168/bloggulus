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

func NewSignin() (*SigninPage, error) {
	// Create the template.
	tmpl := template.New("page")

	// List all required sources.
	sources := []string{
		layout.BaseHTML,
		SigninHTML,
	}

	// Parse each source required to render this page.
	for _, source := range sources {
		_, err := tmpl.Parse(source)
		if err != nil {
			return nil, err
		}
	}

	page := SigninPage{
		tmpl: tmpl,
	}
	return &page, nil
}

func (p *SigninPage) Render(w io.Writer, data SigninData) error {
	return p.tmpl.ExecuteTemplate(w, "page", data)
}
