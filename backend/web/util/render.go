package util

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
)

// Generic function type for rendering data to an io.Writer.
type RenderFunc func(w io.Writer) error

// Render an HTML template based on a given RenderFunc.
func Render(w http.ResponseWriter, r *http.Request, code int, render RenderFunc) {
	// Write the template to the buffer, instead of straight to the http.ResponseWriter.
	// If there's an error, call our serverError() helper and then return.
	var buf bytes.Buffer
	err := render(&buf)
	if err != nil {
		InternalServerErrorResponse(w, r, err)
		return
	}

	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to http.ResponseWriter.
	w.WriteHeader(code)

	// Write the contents of the buffer to the http.ResponseWriter. Note: this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	buf.WriteTo(w)
}

// Render an HTML error template based on a given RenderFunc. This function
// differs from Render in that it doesn't call InternalServerErrorResponse
// if something goes wrong since that could cause an infinite recursion.
func RenderError(w http.ResponseWriter, r *http.Request, code int, render RenderFunc) {
	var buf bytes.Buffer
	err := render(&buf)
	if err != nil {
		slog.Error(err.Error())

		code = http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	w.WriteHeader(code)
	buf.WriteTo(w)
}
