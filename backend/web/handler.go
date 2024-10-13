package web

import (
	"io/fs"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/fetch"
	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
	"github.com/theandrew168/bloggulus/backend/web/util"
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
	conf config.Config,
	repo *repository.Repository,
	find *finder.Finder,
	pageFetcher fetch.PageFetcher,
	syncService *service.SyncService,
) http.Handler {
	mux := http.NewServeMux()

	githubConf := oauth2.Config{
		ClientID:     conf.GithubClientID,
		ClientSecret: conf.GithubClientSecret,
		Scopes:       []string{},
		Endpoint:     github.Endpoint,
	}

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
	mux.Handle("/img/", publicFilesHandler)
	mux.Handle("/js/", publicFilesHandler)

	// The main application routes start here.
	mux.Handle("GET /{$}", HandleIndexPage(find))

	// Authenication routes.
	mux.Handle("GET /login", HandleLogin())
	mux.Handle("GET /github/login", HandleGithubLogin(&githubConf))
	mux.Handle("GET /github/callback", HandleGithubCallback(&githubConf, repo))
	mux.Handle("POST /logout", HandleLogoutForm(repo))

	// Public blog routes.
	mux.Handle("GET /blogs", accountRequired(HandleBlogList(find)))
	mux.Handle("POST /blogs/create", accountRequired(HandleBlogCreateForm(repo, find, syncService)))
	mux.Handle("POST /blogs/{blogID}/follow", accountRequired(HandleBlogFollowForm(repo, find)))
	mux.Handle("POST /blogs/{blogID}/unfollow", accountRequired(HandleBlogUnfollowForm(repo, find)))

	// Public page routes.
	mux.Handle("GET /pages", accountRequired(HandlePageList(repo)))
	mux.Handle("POST /pages/create", accountRequired(HandlePageCreateForm(repo, pageFetcher)))
	mux.Handle("POST /pages/{pageID}/unfollow", accountRequired(HandlePageUnfollowForm(repo)))

	// Private (admin only) blog + post routes.
	mux.Handle("GET /blogs/{blogID}", adminRequired(HandleBlogRead(repo)))
	mux.Handle("POST /blogs/{blogID}/delete", adminRequired(HandleBlogDeleteForm(repo)))
	mux.Handle("GET /blogs/{blogID}/posts/{postID}", adminRequired(HandlePostRead(repo)))
	mux.Handle("POST /blogs/{blogID}/posts/{postID}/delete", adminRequired(HandlePostDeleteForm(repo)))

	// Private (admin only) account routes.
	mux.Handle("GET /accounts", adminRequired(HandleAccountList(repo)))
	mux.Handle("POST /accounts/{accountID}/delete", adminRequired(HandleAccountDeleteForm(repo)))

	mux.HandleFunc("GET /toast", func(w http.ResponseWriter, r *http.Request) {
		cookie := util.NewSessionCookie(util.ToastCookieName, "Toasts are awesome!")
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Requests that don't match any of the above handlers get a 404.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		util.NotFoundResponse(w, r)
	})

	// Apply global middleware to all routes.
	handler := middleware.Use(mux,
		middleware.RecoverPanic(),
		middleware.CompressFiles(),
		middleware.PreventCSRF(),
		middleware.SecureHeaders(),
		middleware.LimitRequestBodySize(),
		middleware.Authenticate(repo),
	)

	return handler
}
