package page

import (
	"html/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/theandrew168/bloggulus/backend/web/partial"
)

func newTemplate(name string, sources []string) *template.Template {
	// Create the template and add helpers.
	tmpl := template.New(name).Funcs(sprig.FuncMap())

	partials := []string{
		partial.ButtonHTML,
		partial.InputHTML,
	}

	// Parse each partial into the template
	for _, source := range partials {
		template.Must(tmpl.Parse(source))
	}

	// Parse each source required to render this page.
	for _, source := range sources {
		template.Must(tmpl.Parse(source))
	}

	return tmpl
}
