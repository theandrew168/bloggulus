package api

import (
	"html/template"
	"net/http"
)

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFS(app.templates, "index.html.tmpl")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
