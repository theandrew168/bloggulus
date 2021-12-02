package web

import (
	"html/template"
	"net/http"
)

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, tmpl string) {
	ts, err := template.ParseFS(app.templates, tmpl)
	if err != nil {
		// skip 3 frames to identify original caller
		app.logger.Output(3, err.Error())
		w.WriteHeader(500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		// skip 3 frames to identify original caller
		app.logger.Output(3, err.Error())
		w.WriteHeader(500)
		return
	}
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, 404, "404.html.tmpl")
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, 405, "405.html.tmpl")
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// skip 2 frames to identify original caller
	app.logger.Output(2, err.Error())
	app.errorResponse(w, r, 500, "500.html.tmpl")
}
