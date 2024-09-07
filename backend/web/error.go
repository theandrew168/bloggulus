package web

import (
	"net/http"

	"github.com/theandrew168/bloggulus/backend/web/page"
)

func HandleNotFoundPage() http.Handler {
	tmpl := page.NewError()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)

		data := page.ErrorData{
			StatusCode: 404,
			StatusText: http.StatusText(404),
			Message:    "This page could not be found.",
		}
		err := tmpl.Render(w, data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
}
