package api

import (
	"fmt"
	"net/http"
)

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := writeJSON(w, status, env, nil)
	if err != nil {
		// skip 3 frames to identify original caller
		app.logger.Output(3, err.Error())
		w.WriteHeader(500)
		return
	}
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, 404, message)
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, 405, message)
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// skip 2 frames to identify original caller
	app.logger.Output(2, err.Error())

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, 500, message)
}
