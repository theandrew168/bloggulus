package page

import (
	"html/template"

	"github.com/Masterminds/sprig/v3"
)

func newTemplate(name string, sources []string) *template.Template {
	// Create the template and add helpers.
	tmpl := template.New(name).Funcs(sprig.FuncMap())

	// Parse each source required to render this page.
	for _, source := range sources {
		template.Must(tmpl.Parse(source))
	}

	return tmpl
}
