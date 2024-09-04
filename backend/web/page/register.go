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

func NewRegister() (*RegisterPage, error) {
	// Create the template.
	tmpl := template.New("page")

	// List all required sources.
	sources := []string{
		layout.BaseHTML,
		RegisterHTML,
	}

	// Parse each source required to render this page.
	for _, source := range sources {
		_, err := tmpl.Parse(source)
		if err != nil {
			return nil, err
		}
	}

	page := RegisterPage{
		tmpl: tmpl,
	}
	return &page, nil
}

func (p *RegisterPage) Render(w io.Writer, data RegisterData) error {
	return p.tmpl.ExecuteTemplate(w, "page", data)
}
