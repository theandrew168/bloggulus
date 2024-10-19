package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed privacy.html
var PrivacyPolicyHTML string

type PrivacyPolicyData struct {
	layout.BaseData
}

type PrivacyPolicyPage struct {
	tmpl *template.Template
}

func NewPrivacyPolicy() *PrivacyPolicyPage {
	sources := []string{
		layout.BaseHTML,
		PrivacyPolicyHTML,
	}

	tmpl := newTemplate("default", sources)
	page := PrivacyPolicyPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *PrivacyPolicyPage) Render(w io.Writer, data PrivacyPolicyData) error {
	return p.tmpl.ExecuteTemplate(w, "default", data)
}
