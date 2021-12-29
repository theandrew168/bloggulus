package web

import (
	"context"
	"html/template"
	"net/http"
	"strconv"

	"github.com/theandrew168/bloggulus/internal/core"
)

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"index.page.tmpl",
		"base.layout.tmpl",
	}

	ts, err := template.ParseFS(app.templates, files...)
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
		ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		// search if requested
		count, err = app.storage.CountSearchPosts(ctx, q)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		posts, err = app.storage.SearchPosts(ctx, q, pageSize, p*pageSize)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		// else just read recent
		count, err = app.storage.CountPosts(ctx)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		posts, err = app.storage.ReadPosts(ctx, pageSize, p*pageSize)
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
		MorePages: (p+1)*pageSize < count,
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
