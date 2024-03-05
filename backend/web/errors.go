package web

import (
	"bytes"
	"html/template"
	"net/http"
)

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, tmpl string) {
	files := []string{
		tmpl,
		"base.layout.tmpl",
	}

	// attempt to parse error template
	ts, err := template.ParseFS(app.templates, files...)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "backend server error", 500)
		return
	}

	// render template to a temp buffer
	var buf bytes.Buffer
	err = ts.Execute(&buf, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "backend server error", 500)
		return
	}

	// write the status and error page
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, 404, "404.page.tmpl")
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, 405, "405.page.tmpl")
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// skip 2 frames to identify original caller
	app.logger.Output(2, err.Error())
	app.errorResponse(w, r, 500, "500.page.tmpl")
}
