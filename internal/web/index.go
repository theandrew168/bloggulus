package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/theandrew168/bloggulus/internal/core"
)

type indexData struct {
	Authed  bool
	Success string
	Error   string
	Posts   []core.Post
}

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFS(app.TemplatesFS, "index.html.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// read the recent posts
	posts, err := app.Post.ReadRecent(r.Context(), 10)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := &indexData{
		Posts:  posts,
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}
