package web

import (
	"io/fs"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
)

// Redirects:
// 303 See Other         - for GETs after POSTs (like a login / register form)
// 302 Found             - all other temporary redirects
// 301 Moved Permanently - permanent redirects

// Route Handler Naming Ideas:
//
// basic page handlers:
// GET - handleIndex
// GET - handleDashboard
//
// basic page w/ form handlers:
// GET  - handleLogin
// POST - handleLoginForm
//
// CRUD handlers:
// C POST - handleFooCreate[Form]
// R GET  - handleFooRead
// U POST - handleFooUpdate[Form]
// D POST - handleFooDelete[Form]
// L GET  - handleFooList

func Handler(
	public fs.FS,
	repo *repository.Repository,
	find *finder.Finder,
	syncService *service.SyncService,
) http.Handler {
	mux := http.NewServeMux()

	accountRequired := middleware.AccountRequired()
	adminRequired := middleware.Chain(accountRequired, middleware.AdminRequired())

	// Host prometheus metrics on "/metrics".
	mux.Handle("GET /metrics", promhttp.Handler())

	// Basic health check endpoint for verifying that the app is running.
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	})

	// Compress and serve the embedded public (static) files.
	publicFiles, _ := fs.Sub(public, "public")
	publicFilesHandler := http.FileServer(http.FS(publicFiles))

	// Serve public (static) files from the embedded FS.
	mux.Handle("/favicon.ico", publicFilesHandler)
	mux.Handle("/robots.txt", publicFilesHandler)
	mux.Handle("/css/", publicFilesHandler)
	mux.Handle("/js/", publicFilesHandler)

	// The main application routes start here.
	mux.Handle("GET /{$}", HandleIndexPage(find))

	// Authenication routes.
	mux.Handle("GET /register", HandleRegisterPage())
	mux.Handle("POST /register", HandleRegisterForm(repo))
	mux.Handle("GET /signin", HandleSigninPage())
	mux.Handle("POST /signin", HandleSigninForm(repo))
	mux.Handle("POST /signout", HandleSignoutForm(repo))

	// Blog and post routes.
	mux.Handle("GET /blogs", accountRequired(HandleBlogList(find)))
	mux.Handle("POST /blogs/create", accountRequired(HandleBlogCreateForm(repo, find, syncService)))
	mux.Handle("POST /blogs/follow", accountRequired(HandleBlogFollowForm(repo, find)))
	mux.Handle("POST /blogs/unfollow", accountRequired(HandleBlogUnfollowForm(repo, find)))
	mux.Handle("GET /blogs/{blogID}", accountRequired(HandleBlogRead(repo)))
	mux.Handle("POST /blogs/{blogID}/delete", adminRequired(HandleBlogDeleteForm(repo)))

	// Requests that don't match any of the above handlers get a 404.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("<h1>Bloggulus - Not Found</h1>"))
	})

	// Apply global middleware to all routes.
	handler := middleware.Use(mux,
		middleware.RecoverPanic(),
		middleware.CompressFiles(),
		middleware.SecureHeaders(),
		middleware.LimitRequestBodySize(),
		middleware.Authenticate(repo),
	)

	return handler
}
