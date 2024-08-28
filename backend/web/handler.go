package web

import (
	"io/fs"
	"net/http"

	"github.com/klauspost/compress/gzhttp"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
	"github.com/theandrew168/bloggulus/backend/web/page"
)

func Handler(
	public fs.FS,
	store *storage.Storage,
	syncService *service.SyncService,
) http.Handler {
	mux := http.NewServeMux()

	// metrics
	mux.Handle("GET /metrics", promhttp.Handler())

	// basic health check
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	})

	// compress public files
	publicFiles, _ := fs.Sub(public, "public")
	publicHandler := gzhttp.GzipHandler(http.FileServer(http.FS(publicFiles)))

	// serve non-index static files from the public FS
	mux.Handle("/favicon.ico", publicHandler)
	mux.Handle("/robots.txt", publicHandler)
	mux.Handle("/css/", publicHandler)
	mux.Handle("/fonts/", publicHandler)

	mux.Handle("/{$}", page.HandleIndex(store))

	return middleware.Use(mux, middleware.RecoverPanic())
}
