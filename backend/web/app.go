package web

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/klauspost/compress/gzhttp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
	metricsWrapper "github.com/slok/go-http-metrics/middleware/std"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/web/api"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
)

type Application struct {
	logger   *log.Logger
	storage  storage.Storage
	frontend fs.FS
}

func NewApplication(logger *log.Logger, storage storage.Storage, frontend fs.FS) *Application {
	app := Application{
		logger:   logger,
		storage:  storage,
		frontend: frontend,
	}
	return &app
}

func (app *Application) Router() http.Handler {
	mmw := metricsMiddleware.New(metricsMiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	mux := flow.New()
	mux.Use(middleware.RecoverPanic)

	// metrics
	mux.Handle("/metrics", promhttp.Handler(), "GET")

	// basic health check
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	}, "GET")

	// backend - rest api
	apiApp := api.NewApplication(app.logger, app.storage)
	mux.Handle("/api/v1/...", metricsWrapper.Handler("/api/v1", mmw, http.StripPrefix("/api/v1", apiApp.Router())))
	mux.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/", http.StatusMovedPermanently)
	})

	// frontend - svelte
	frontendHandler := gzhttp.GzipHandler(http.FileServer(http.FS(app.frontend)))

	mux.Handle("/", frontendHandler)
	mux.Handle("/index.html", frontendHandler)
	mux.Handle("/robots.txt", frontendHandler)
	mux.Handle("/favicon.png", frontendHandler)
	mux.Handle("/openapi.yaml", frontendHandler)
	mux.Handle("/_app/...", frontendHandler)

	// all other routes should return the index page
	// so that the frontend router can take over
	mux.Handle("/...", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		index, err := fs.ReadFile(app.frontend, "index.html")
		if err != nil {
			panic(err)
		}

		w.Write(index)
	}))

	return mux
}
