package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed jobs.html
var JobsHTML string

type JobsData struct {
	layout.BaseData

	Jobs []model.DeadJob
}

type JobsPage struct {
	tmpl *template.Template
}

func NewJobs() *JobsPage {
	sources := []string{
		layout.BaseHTML,
		JobsHTML,
	}

	tmpl := newTemplate("default", sources)
	page := JobsPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *JobsPage) Render(w io.Writer, data JobsData) error {
	data.Stylesheets = []string{
		"/css/jobs.css",
	}
	return p.tmpl.ExecuteTemplate(w, "default", data)
}
