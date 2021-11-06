package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/theandrew168/bloggulus/internal/core"
)

type Application struct {
	Logger *log.Logger

	Blog core.BlogStorage
	Post core.PostStorage
}

func (app *Application) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello API!"))
	})
	return r
}
