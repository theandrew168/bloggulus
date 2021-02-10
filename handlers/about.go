package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func (app *Application) HandleAbout(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/about.html.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
}
