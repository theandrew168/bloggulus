package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/query"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed blogs.html
var BlogsHTML string

type BlogsData struct {
	layout.BaseData

	Blogs []BlogsBlogData
}

// Since this data might be re-rendered per-row via HTMX, we have to include
// the layout.BaseData to ensure CSRF prevention still works.
type BlogsBlogData struct {
	layout.BaseData
	query.BlogForAccount
}

type BlogsPage struct {
	tmpl *template.Template
}

func NewBlogs() *BlogsPage {
	sources := []string{
		layout.BaseHTML,
		BlogsHTML,
	}

	tmpl := newTemplate("page", sources)
	page := BlogsPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *BlogsPage) Render(w io.Writer, data BlogsData) error {
	return p.tmpl.ExecuteTemplate(w, "page", data)
}

func (p *BlogsPage) RenderBlog(w io.Writer, data BlogsBlogData) error {
	return p.tmpl.ExecuteTemplate(w, "blog", data)
}
