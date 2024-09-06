package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed blogs.html
var BlogsHTML string

type BlogsData struct {
	layout.BaseData

	Account *model.Account
	Blogs   []finder.BlogForAccount
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

func (p *BlogsPage) RenderBlog(w io.Writer, data finder.BlogForAccount) error {
	return p.tmpl.ExecuteTemplate(w, "blog", data)
}

func (p *BlogsPage) RenderBlogs(w io.Writer, data BlogsData) error {
	return p.tmpl.ExecuteTemplate(w, "blogs", data)
}
