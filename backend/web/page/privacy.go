package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed privacy.html
var PrivacyHTML string

type PrivacyData struct {
	layout.BaseData
}

type PrivacyPage struct {
	tmpl *template.Template
}

func NewPrivacy() *PrivacyPage {
	sources := []string{
		layout.BaseHTML,
		PrivacyHTML,
	}

	tmpl := newTemplate("default", sources)
	page := PrivacyPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *PrivacyPage) Render(w io.Writer, data PrivacyData) error {
	return p.tmpl.ExecuteTemplate(w, "default", data)
}
