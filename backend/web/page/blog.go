package page

import (
	_ "embed"
	"html/template"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed blog.html
var BlogHTML string

// TODO: Come up with a better name for this? BlogPagePost? BlogPost?
type BlogDataPost struct {
	ID          uuid.UUID
	BlogID      uuid.UUID
	Title       string
	PublishedAt time.Time
}

type BlogData struct {
	layout.BaseData

	ID       uuid.UUID
	Title    string
	SiteURL  string
	FeedURL  string
	SyncedAt time.Time
	Posts    []BlogDataPost
}

type BlogPage struct {
	tmpl *template.Template
}

func NewBlog() *BlogPage {
	sources := []string{
		layout.BaseHTML,
		BlogHTML,
	}

	tmpl := newTemplate("page", sources)
	page := BlogPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *BlogPage) Render(w io.Writer, data BlogData) error {
	return p.tmpl.ExecuteTemplate(w, "page", data)
}
