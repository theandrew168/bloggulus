package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/theandrew168/bloggulus/models"
)

type indexData struct {
	Authed bool
	Posts  []*models.SourcedPost
}

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/index.html.tmpl", "templates/base.html.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	accountID, err := app.CheckSessionAccount(w, r)
	if err != nil {
		if err != ErrNoSession {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	authed := err == nil

	var posts []*models.SourcedPost
	if authed {
		posts, err = app.SourcedPost.ReadRecentForUser(r.Context(), accountID, 10)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		posts, err = app.SourcedPost.ReadRecent(r.Context(), 10)
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
