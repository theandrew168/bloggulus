package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/bloggulus/internal/api"
	"github.com/theandrew168/bloggulus/internal/static"
	"github.com/theandrew168/bloggulus/internal/storage"
	"github.com/theandrew168/bloggulus/internal/web"
)

func New(logger *log.Logger, storage *storage.Storage) http.Handler {
	mux := chi.NewRouter()

	// handle top-level special cases
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	})
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(static.Favicon)
	})
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(static.Robots)
	})

	// static files app
	staticApp := static.NewApplication()
	mux.Mount("/static", http.StripPrefix("/static", staticApp.Router()))

	// rest api app
	apiApp := api.NewApplication(logger, storage)
	mux.Mount("/api/v1", http.StripPrefix("/api/v1", apiApp.Router()))

	// primary web app (last due to being a top-level catch-all)
	webApp := web.NewApplication(logger, storage)
	mux.Mount("/", webApp.Router())

	return mux
}
