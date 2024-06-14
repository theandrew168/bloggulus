package web

import (
	"io/fs"
	"net/http"

	"github.com/klauspost/compress/gzhttp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
	metricsWrapper "github.com/slok/go-http-metrics/middleware/std"

	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/api"
	"github.com/theandrew168/bloggulus/backend/web/api/middleware"
)

type Application struct {
	frontend fs.FS
	store    *storage.Storage
}

func NewApplication(
	frontend fs.FS,
	store *storage.Storage,
) *Application {
	app := Application{
		frontend: frontend,
		store:    store,
	}
	return &app
}

func (app *Application) Router() http.Handler {
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
	apiApp := api.NewApplication(app.store)
	mux.Handle("/api/v1/", metricsWrapper.Handler("/api/v1", mmw, http.StripPrefix("/api/v1", apiApp.Router())))

	// frontend - svelte
	frontendHandler := gzhttp.GzipHandler(http.FileServer(http.FS(app.frontend)))

	// serve non-index static files from the frontend FS
	mux.Handle("/robots.txt", frontendHandler)
	mux.Handle("/favicon.png", frontendHandler)
	mux.Handle("/openapi.yaml", frontendHandler)
	mux.Handle("/fonts/", frontendHandler)
	mux.Handle("/_app/", frontendHandler)

	// all other routes should return the index page so that the frontend router can take over
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		index, err := fs.ReadFile(app.frontend, "index.html")
		if err != nil {
			panic(err)
		}

		w.Write(index)
	}))

	return middleware.Adapt(mux, middleware.RecoverPanic())
}
