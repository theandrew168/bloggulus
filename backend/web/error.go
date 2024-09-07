package web

import (
	"io"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/web/page"
)

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	tmpl := page.NewError()

	code := 404
	data := page.ErrorData{
		StatusCode: code,
		StatusText: http.StatusText(code),
		Message:    "Sorry, this page could not be found.",
	}

	RenderError(w, r, code, func(w io.Writer) error {
		return tmpl.Render(w, data)
	})
}

func InternalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	tmpl := page.NewError()

	code := 500
	data := page.ErrorData{
		StatusCode: code,
		StatusText: http.StatusText(code),
		Message:    "Sorry, something went wrong.",
	}

	RenderError(w, r, code, func(w io.Writer) error {
		return tmpl.Render(w, data)
	})
}

func HandleNotFoundPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		NotFoundResponse(w, r)
	})
}
