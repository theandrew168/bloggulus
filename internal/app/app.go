package app

import (
	"log"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
	metricsWrapper "github.com/slok/go-http-metrics/middleware/std"

	"github.com/theandrew168/bloggulus/internal/api"
	"github.com/theandrew168/bloggulus/internal/static"
	"github.com/theandrew168/bloggulus/internal/storage"
	"github.com/theandrew168/bloggulus/internal/web"
)

func New(logger *log.Logger, storage *storage.Storage) http.Handler {
	mmw := metricsMiddleware.New(metricsMiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	mux := flow.New()

	// handle top-level special cases
	mux.Handle("/metrics", promhttp.Handler(), "GET")
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	}, "GET")
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(static.Favicon)
	}, "GET")
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(static.Robots)
	}, "GET")

	// static files app
	staticApp := static.NewApplication()
	mux.Handle("/static/...", http.StripPrefix("/static", staticApp.Router()), "GET")
	mux.HandleFunc("/static", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/", http.StatusMovedPermanently)
	})

	// rest api app
	apiApp := api.NewApplication(logger, storage)
	mux.Handle("/api/v1/...", metricsWrapper.Handler("/api/v1", mmw, http.StripPrefix("/api/v1", apiApp.Router())))
	mux.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/", http.StatusMovedPermanently)
	})

	// primary web app (last due to being a top-level catch-all)
	webApp := web.NewApplication(logger, storage)
	mux.Handle("/...", metricsWrapper.Handler("/", mmw, webApp.Router()), "GET")

	return mux
}
