package web

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/bloggulus/internal/core"
)

type Application struct {
	StaticFS    fs.FS
	TemplatesFS fs.FS

	Blog core.BlogStorage
	Post core.PostStorage
}

func (app *Application) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", app.HandleIndex)
	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.FS(app.StaticFS))))
	return r
}
