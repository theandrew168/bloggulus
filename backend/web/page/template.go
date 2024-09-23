package page

import (
	"html/template"
)

func newTemplate(name string, sources []string) *template.Template {
	// Create the template and add helpers.
	tmpl := template.New(name)

	// Parse each source required to render this page.
	for _, source := range sources {
		template.Must(tmpl.Parse(source))
	}

	return tmpl
}
