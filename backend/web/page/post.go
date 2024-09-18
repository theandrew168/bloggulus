package page

import (
	_ "embed"
	"html/template"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed post.html
var PostHTML string

type PostData struct {
	layout.BaseData

	ID          uuid.UUID
	BlogID      uuid.UUID
	Title       string
	URL         string
	PublishedAt time.Time
}

type PostPage struct {
	tmpl *template.Template
}

func NewPost() *PostPage {
	sources := []string{
		layout.BaseHTML,
		PostHTML,
	}

	tmpl := newTemplate("page", sources)
	page := PostPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *PostPage) Render(w io.Writer, data PostData) error {
	return p.tmpl.ExecuteTemplate(w, "page", data)
}
