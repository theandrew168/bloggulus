package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed post.html
var PostHTML string

type PostData struct {
	layout.BaseData

	Post *model.Post
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
	return p.tmpl.ExecuteTemplate(w, "default", data)
}
