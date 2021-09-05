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
	Posts   []*core.Post
}

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/index.html.tmpl", "templates/base.html.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	account, err := app.CheckAccount(w, r)
	if err != nil {
		if err != ErrNoSession {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	authed := err == nil

	var posts []*core.Post
	if authed {
		// read the recent posts that the user follows
		posts, err = app.Post.ReadRecentForUser(r.Context(), account.AccountID, 10)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		// read the recent posts
		posts, err = app.Post.ReadRecent(r.Context(), 10)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	data := &indexData{
		Authed: authed,
		Posts:  posts,
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}