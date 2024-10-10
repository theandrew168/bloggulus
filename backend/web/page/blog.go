package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed blog.html
var BlogHTML string

type BlogData struct {
	layout.BaseData

	Blog  *model.Blog
	Posts []*model.Post
}

type BlogPage struct {
	tmpl *template.Template
}

func NewBlog() *BlogPage {
	sources := []string{
		layout.BaseHTML,
		BlogHTML,
	}

	tmpl := newTemplate("default", sources)
	page := BlogPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *BlogPage) Render(w io.Writer, data BlogData) error {
	return p.tmpl.ExecuteTemplate(w, "default", data)
}
