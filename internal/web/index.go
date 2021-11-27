package web

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/theandrew168/bloggulus/internal/core"
)

const (
	PageSize = 15
)

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFS(app.templates, "index.html.tmpl")
	if err != nil {
		http.Error(w, "Internal server error", 500)
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
		count, err = app.storage.PostCountSearch(r.Context(), q)
		if err != nil {
			http.Error(w, "Internal server error", 500)
			return
		}

		posts, err = app.storage.PostReadSearch(r.Context(), q, PageSize, p*PageSize)
		if err != nil {
			http.Error(w, "Internal server error", 500)
			return
		}

		// else just read recent
	} else {
		count, err = app.storage.PostCountRecent(r.Context())
		if err != nil {
			http.Error(w, "Internal server error", 500)
			return
		}

		posts, err = app.storage.PostReadRecent(r.Context(), PageSize, p*PageSize)
		if err != nil {
			http.Error(w, "Internal server error", 500)
			return
		}
	}

	// limit each post to 3 tags
	for i := 0; i < len(posts); i++ {
		if len(posts[i].Tags) > 3 {
			posts[i].Tags = posts[i].Tags[:3]
		}
	}

	data := struct {
		MorePages bool
		NextPage  int
		Search    string
		Posts     []core.Post
	}{
		MorePages: (p+1)*PageSize < count,
		NextPage:  p + 1,
		Search:    q,
		Posts:     posts,
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.logger.Println(err)
		return
	}
}
