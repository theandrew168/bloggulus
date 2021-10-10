package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/theandrew168/bloggulus/internal/core"
)

const (
	PageSize = 15
)

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFS(app.TemplatesFS, "index.html.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// check page param
	p, err := strconv.Atoi(r.URL.Query().Get("p"))
	if err != nil {
		p = 0
	}

	// check search param
	q := r.URL.Query().Get("q")

	var count int
	var posts []core.Post

	// search if requested
	if q != "" {
		count, err = app.Post.CountSearch(r.Context(), q)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		posts, err = app.Post.ReadSearch(r.Context(), q, PageSize, p * PageSize)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

	// else just read recent
	} else {
		count, err = app.Post.CountRecent(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		posts, err = app.Post.ReadRecent(r.Context(), PageSize, p * PageSize)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	data := struct {
		MorePages bool
		NextPage  int
		Search    string
		Posts     []core.Post
	}{
		MorePages: (p + 1) * PageSize < count,
		NextPage:  p + 1,
		Search:    q,
		Posts:     posts,
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}
