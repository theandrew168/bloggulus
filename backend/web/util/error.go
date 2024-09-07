package util

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/web/page"
)

// Handle errors that arise from creating a new row in the database.
func CreateErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, postgres.ErrConflict):
		BadRequestResponse(w, r)
	default:
		InternalServerErrorResponse(w, r, err)
	}
}

// Handle errors that arise from reading a single row from the database.
func ReadErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, postgres.ErrNotFound):
		NotFoundResponse(w, r)
	default:
		InternalServerErrorResponse(w, r, err)
	}
}

// Handle errors that arise from reading many (zero or more) rows from the database.
func ListErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	InternalServerErrorResponse(w, r, err)
}

// Handle errors that arise from updating a single row in the database.
func UpdateErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, postgres.ErrNotFound):
		NotFoundResponse(w, r)
	case errors.Is(err, postgres.ErrConflict):
		BadRequestResponse(w, r)
	default:
		InternalServerErrorResponse(w, r, err)
	}
}

// Handle errors that arise from deleting a single row from the database.
func DeleteErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, postgres.ErrNotFound):
		NotFoundResponse(w, r)
	default:
		InternalServerErrorResponse(w, r, err)
	}
}

// Render a 400 Bad Request error page.
func BadRequestResponse(w http.ResponseWriter, r *http.Request) {
	tmpl := page.NewError()
	code := http.StatusBadRequest
	data := page.ErrorData{
		StatusCode: code,
		StatusText: http.StatusText(code),
		Message:    "Sorry, we cannot understand what you sent.",
	}

	RenderError(w, r, code, func(w io.Writer) error {
		return tmpl.Render(w, data)
	})
}

// Render a 403 Forbidden error page.
func ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	tmpl := page.NewError()
	code := http.StatusForbidden
	data := page.ErrorData{
		StatusCode: code,
		StatusText: http.StatusText(code),
		Message:    "Sorry, you are not allowed to access this page.",
	}

	RenderError(w, r, code, func(w io.Writer) error {
		return tmpl.Render(w, data)
	})
}

// Render a 404 Not Found error page.
func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	tmpl := page.NewError()
	code := http.StatusNotFound
	data := page.ErrorData{
		StatusCode: code,
		StatusText: http.StatusText(code),
		Message:    "Sorry, this page could not be found.",
	}

	RenderError(w, r, code, func(w io.Writer) error {
		return tmpl.Render(w, data)
	})
}

// Render a 500 Internal Server Error page.
func InternalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error(err.Error(), "url", r.URL.String())

	tmpl := page.NewError()
	code := http.StatusInternalServerError
	data := page.ErrorData{
		StatusCode: code,
		StatusText: http.StatusText(code),
		Message:    "Sorry, something went wrong.",
	}

	RenderError(w, r, code, func(w io.Writer) error {
		return tmpl.Render(w, data)
	})
}
