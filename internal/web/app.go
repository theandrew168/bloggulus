package web

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/theandrew168/bloggulus/internal/core"
)

type Application struct {
	TemplatesFS fs.FS

	Blog core.BlogStorage
	Post core.PostStorage
}

func (app *Application) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", app.HandleIndex)
	return r
}
