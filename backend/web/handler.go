package web

import (
	"io/fs"
	"net/http"

	"github.com/klauspost/compress/gzhttp"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
)

func Handler(
	public fs.FS,
	repo *repository.Repository,
	find *finder.Finder,
	syncService *service.SyncService,
) http.Handler {
	mux := http.NewServeMux()

	accountRequired := middleware.AccountRequired()
	// adminRequired := middleware.Chain(accountRequired, middleware.AdminRequired())

	// Host prometheus metrics on "/metrics".
	mux.Handle("GET /metrics", promhttp.Handler())

	// Basic health check endpoint for verifying that the app is running.
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	})

	// Compress and serve the embedded public (static) files.
	publicFiles, _ := fs.Sub(public, "public")
	publicFilesHandler := gzhttp.GzipHandler(http.FileServer(http.FS(publicFiles)))

	// Serve public (static) files from the embedded FS.
	mux.Handle("/favicon.ico", publicFilesHandler)
	mux.Handle("/robots.txt", publicFilesHandler)
	mux.Handle("/css/", publicFilesHandler)
	mux.Handle("/js/", publicFilesHandler)

	// The main application routes start here.
	mux.Handle("GET /{$}", HandleIndexPage(find))

	mux.Handle("GET /register", HandleRegisterPage())
	mux.Handle("POST /register", HandleRegisterForm(repo))
	mux.Handle("GET /signin", HandleSigninPage())
	mux.Handle("POST /signin", HandleSigninForm(repo))
	mux.Handle("POST /signout", HandleSignoutForm(repo))

	mux.Handle("GET /blogs", accountRequired(HandleBlogsPage(find)))
	mux.Handle("POST /blogs", accountRequired(HandleBlogsForm(repo, find, syncService)))

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
		middleware.Authenticate(repo),
	)
}
