package api

import (
	"net/http"
)

func (app *Application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
		return
	}
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, 500, message)
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, 404, message)
}
