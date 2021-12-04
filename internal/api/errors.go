package api

import (
	"net/http"
)

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := writeJSON(w, status, env)
	if err != nil {
		app.logger.Println(err)
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
	// skip 2 frames to identify original caller
	app.logger.Output(2, err.Error())

	message := "internal server error"
	app.errorResponse(w, r, 500, message)
}
