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
		app.serverErrorResponse(w, r, err)
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

	if q != "" {
		// search if requested
		count, err = app.storage.CountSearchPosts(r.Context(), q)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		posts, err = app.storage.SearchPosts(r.Context(), q, PageSize, p*PageSize)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else {
		// else just read recent
		count, err = app.storage.CountPosts(r.Context())
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		posts, err = app.storage.ReadPosts(r.Context(), PageSize, p*PageSize)
		if err != nil {
			app.serverErrorResponse(w, r, err)
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
		app.serverErrorResponse(w, r, err)
		return
	}
}
