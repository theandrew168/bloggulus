package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/query"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed post.html
var PostHTML string

type PostData struct {
	layout.BaseData

	Post query.PostDetails
}

type PostPage struct {
	tmpl *template.Template
}

func NewPost() *PostPage {
	sources := []string{
		layout.BaseHTML,
		PostHTML,
	}

	tmpl := newTemplate("default", sources)
	page := PostPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *PostPage) Render(w io.Writer, data PostData) error {
	data.Stylesheets = []string{
		"/css/details.css",
	}
	return p.tmpl.ExecuteTemplate(w, "default", data)
}
