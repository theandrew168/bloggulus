package web

import (
	"html/template"
	"log"
	"net/http"
)

type aboutData struct {
	Authed bool
}

func (app *Application) HandleAbout(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/about.html.tmpl", "templates/base.html.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = app.CheckSessionAccount(w, r)
	if err != nil {
		if err != ErrNoSession {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	authed := err == nil

	data := &aboutData{
		Authed: authed,
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}
