package api

import (
	"log/slog"
	"net/http"
)

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	type response struct {
		Error any `json:"error"`
	}
	resp := response{
		Error: message,
	}

	err := writeJSON(w, status, resp, nil)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(500)
		return
	}
}

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, 400, errors)
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "not found"
	app.errorResponse(w, r, 404, message)
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := "method not allowed"
	app.errorResponse(w, r, 405, message)
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error(err.Error())

	message := "backend server error"
	app.errorResponse(w, r, 500, message)
}
