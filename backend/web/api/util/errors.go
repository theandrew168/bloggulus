package util

import (
	"log/slog"
	"net/http"
	"strings"
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

func BadRequestResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusBadRequest
	text := http.StatusText(code)
	ErrorResponse(w, r, code, strings.ToLower((text)))
}

func UnauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusUnauthorized
	text := http.StatusText(code)
	ErrorResponse(w, r, code, strings.ToLower((text)))
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusNotFound
	text := http.StatusText(code)
	ErrorResponse(w, r, code, strings.ToLower((text)))
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusMethodNotAllowed
	text := http.StatusText(code)
	ErrorResponse(w, r, code, strings.ToLower((text)))
}

func ConflictResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusConflict
	text := http.StatusText(code)
	ErrorResponse(w, r, code, strings.ToLower((text)))
}

// Note that the errors parameter here has the type map[string]string,
// which is exactly the same as the errors map contained in our Validator type.
func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	code := http.StatusUnprocessableEntity
	ErrorResponse(w, r, code, errors)
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error(err.Error())

	code := http.StatusInternalServerError
	text := http.StatusText(code)
	ErrorResponse(w, r, code, strings.ToLower(text))
}
