package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed blogs.html
var BlogsHTML string

type BlogsData struct {
	layout.BaseData

	Blogs []finder.BlogForAccount
}

type BlogsPage struct {
	tmpl *template.Template
}

func NewBlogs() (*BlogsPage, error) {
	// Create the template.
	tmpl := template.New("page")

	// List all required sources.
	sources := []string{
		layout.BaseHTML,
		BlogsHTML,
	}

	// Parse each source required to render this page.
	for _, source := range sources {
		_, err := tmpl.Parse(source)
		if err != nil {
			return nil, err
		}
	}

	page := BlogsPage{
		tmpl: tmpl,
	}
	return &page, nil
}

func (p *BlogsPage) Render(w io.Writer, data BlogsData) error {
	return p.tmpl.ExecuteTemplate(w, "page", data)
}

func (p *BlogsPage) RenderBlog(w io.Writer, data finder.BlogForAccount) error {
	return p.tmpl.ExecuteTemplate(w, "blog", data)
}

func (p *BlogsPage) RenderBlogs(w io.Writer, data BlogsData) error {
	return p.tmpl.ExecuteTemplate(w, "blogs", data)
}
