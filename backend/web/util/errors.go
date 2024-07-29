package util

import (
	"log/slog"
	"net/http"
	"strings"
)

type Error struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

type Errors struct {
	Errors []Error `json:"errors"`
}

func NewErrors() *Errors {
	e := Errors{}
	return &e
}

func (e *Errors) Add(message string) {
	e.AddField(message, "")
}

func (e *Errors) AddField(message, field string) {
	e.Errors = append(e.Errors, Error{Message: message, Field: field})
}

func (e *Errors) Check(ok bool, message string) {
	if !ok {
		e.Add(message)
	}
}

func (e *Errors) CheckField(ok bool, message, field string) {
	if !ok {
		e.AddField(message, field)
	}
}

func (e *Errors) Valid() bool {
	return len(e.Errors) == 0
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, errors *Errors) {
	err := WriteJSON(w, status, errors, nil)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusBadRequest
	text := http.StatusText(code)
	errors := NewErrors()
	errors.Add(strings.ToLower(text))
	ErrorResponse(w, r, code, errors)
}

func UnauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusUnauthorized
	text := http.StatusText(code)
	errors := NewErrors()
	errors.Add(strings.ToLower(text))
	ErrorResponse(w, r, code, errors)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusNotFound
	text := http.StatusText(code)
	errors := NewErrors()
	errors.Add(strings.ToLower(text))
	ErrorResponse(w, r, code, errors)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusMethodNotAllowed
	text := http.StatusText(code)
	errors := NewErrors()
	errors.Add(strings.ToLower(text))
	ErrorResponse(w, r, code, errors)
}

func ConflictResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusConflict
	text := http.StatusText(code)
	errors := NewErrors()
	errors.Add(strings.ToLower(text))
	ErrorResponse(w, r, code, errors)
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors *Errors) {
	code := http.StatusUnprocessableEntity
	ErrorResponse(w, r, code, errors)
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error(err.Error())

	code := http.StatusInternalServerError
	text := http.StatusText(code)
	errors := NewErrors()
	errors.Add(strings.ToLower(text))
	ErrorResponse(w, r, code, errors)
}
