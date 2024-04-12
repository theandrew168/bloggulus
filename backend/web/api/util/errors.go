package util

import (
	"log/slog"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	type response struct {
		Error any `json:"error"`
	}
	resp := response{
		Error: message,
	}

	err := WriteJSON(w, status, resp, nil)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(500)
		return
	}
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	ErrorResponse(w, r, 400, errors)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "not found"
	ErrorResponse(w, r, 404, message)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := "method not allowed"
	ErrorResponse(w, r, 405, message)
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error(err.Error())

	message := "backend server error"
	ErrorResponse(w, r, 500, message)
}
