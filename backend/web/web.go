package web

import (
	"io/fs"
	"net/http"

	"github.com/klauspost/compress/gzhttp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
	metricsWrapper "github.com/slok/go-http-metrics/middleware/std"

	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/api"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
)

func Handler(
	frontend fs.FS,
	store *storage.Storage,
	syncService *service.SyncService,
) http.Handler {
	mmw := metricsMiddleware.New(metricsMiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	mux := http.NewServeMux()

	// metrics
	mux.Handle("GET /metrics", promhttp.Handler())

	// basic health check
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	})

	// backend - rest api
	apiHandler := api.Handler(store, syncService)
	mux.Handle("/api/v1/", metricsWrapper.Handler("/api/v1", mmw, http.StripPrefix("/api/v1", apiHandler)))

	// frontend - svelte
	frontendHandler := gzhttp.GzipHandler(http.FileServer(http.FS(frontend)))

	// serve non-index static files from the frontend FS
	mux.Handle("/favicon.png", frontendHandler)
	mux.Handle("/openapi.yaml", frontendHandler)
	mux.Handle("/robots.txt", frontendHandler)
	mux.Handle("/assets/", frontendHandler)
	mux.Handle("/css/", frontendHandler)
	mux.Handle("/fonts/", frontendHandler)

	// all other routes should return the index page so that the frontend router can take over
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		index, err := fs.ReadFile(frontend, "index.html")
		if err != nil {
			panic(err)
		}

		w.Write(index)
	}))

	return middleware.Use(mux, middleware.RecoverPanic())
}
