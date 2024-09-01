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

	// Host prometheus metrics on "/metrics".
	mux.Handle("GET /metrics", promhttp.Handler())

	// Basic health check endpoint for verifying that the app is running.
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	})

	// Compress and serve the embedded public (static) files.
	publicFiles, _ := fs.Sub(public, "public")
	publicHandler := gzhttp.GzipHandler(http.FileServer(http.FS(publicFiles)))

	// Serve public (static) files from the embedded FS.
	mux.Handle("/favicon.ico", publicHandler)
	mux.Handle("/robots.txt", publicHandler)
	mux.Handle("/css/", publicHandler)

	// The main application routes start here.
	mux.Handle("/{$}", page.HandleIndex(store))
	mux.Handle("GET /register", page.HandleRegister())
	mux.Handle("POST /register", page.HandleRegisterForm(store))
	mux.Handle("GET /signin", page.HandleSignin())
	mux.Handle("POST /signin", page.HandleSigninForm(store))
	mux.Handle("POST /signout", page.HandleSignoutForm(store))

	// Requests that don't match any of the above handlers get a 404.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("<h1>Bloggulus - Not Found</h1>"))
	})

	// Apply global middleware to all routes.
	return middleware.Use(mux,
		middleware.RecoverPanic(),
		middleware.SecureHeaders(),
		middleware.LimitRequestBodySize(),
	)
}
