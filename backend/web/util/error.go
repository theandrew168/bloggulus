package util

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/web/ui"
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
	code := http.StatusBadRequest
	page := ui.ErrorPage(ui.ErrorPageData{
		PageLayoutData: GetPageLayoutData(r, w),

		StatusCode: code,
		StatusText: http.StatusText(code),
		Message:    "Sorry, we cannot understand what you sent.",
	})

	RenderError(w, r, code, func(w io.Writer) error {
		return page.Render(w)
	})
}

// Render a 403 Forbidden error page.
func ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusForbidden
	page := ui.ErrorPage(ui.ErrorPageData{
		PageLayoutData: GetPageLayoutData(r, w),

		StatusCode: code,
		StatusText: http.StatusText(code),
		Message:    "Sorry, you are not allowed to access this page.",
	})

	RenderError(w, r, code, func(w io.Writer) error {
		return page.Render(w)
	})
}

// Render a 404 Not Found error page.
func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusNotFound
	page := ui.ErrorPage(ui.ErrorPageData{
		PageLayoutData: GetPageLayoutData(r, w),

		StatusCode: code,
		StatusText: http.StatusText(code),
		Message:    "Sorry, this page could not be found.",
	})

	RenderError(w, r, code, func(w io.Writer) error {
		return page.Render(w)
	})
}

// Render a 500 Internal Server Error page.
func InternalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("internal server error",
		"error", err.Error(),
		"url", r.URL.String(),
	)

	code := http.StatusInternalServerError
	page := ui.ErrorPage(ui.ErrorPageData{
		PageLayoutData: GetPageLayoutData(r, w),

		StatusCode: code,
		StatusText: http.StatusText(code),
		Message:    "Sorry, something went wrong.",
	})

	RenderError(w, r, code, func(w io.Writer) error {
		return page.Render(w)
	})
}
