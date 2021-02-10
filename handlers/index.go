package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/theandrew168/bloggulus/models"
)

type IndexData struct {
	Posts []*models.SourcedPost
}

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/index.html.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	posts, err := app.SourcedPost.ReadRecent(r.Context(), 20)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := &IndexData{
		Posts: posts,
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}
