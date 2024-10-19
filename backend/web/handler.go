package web

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

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
		Endpoint:     github.Endpoint,
		ClientID:     conf.GithubClientID,
		ClientSecret: conf.GithubClientSecret,
		RedirectURL:  conf.GithubRedirectURI,
		Scopes:       []string{},
	}
	googleConf := oauth2.Config{
		Endpoint:     google.Endpoint,
		ClientID:     conf.GoogleClientID,
		ClientSecret: conf.GoogleClientSecret,
		RedirectURL:  conf.GoogleRedirectURI,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
	}

	requireAccount := middleware.RequireAccount()
	requireAdmin := middleware.Chain(requireAccount, middleware.RequireAdmin())

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

	// Documentation routes.
	mux.Handle("GET /privacy", HandlePrivacyPolicy())

	// Check if the debug auth method should be enabled.
	enableDebugAuth := os.Getenv("ENABLE_DEBUG_AUTH") != ""
	if enableDebugAuth {
		mux.Handle("POST /debug/signin", HandleDebugSignIn(repo))
	}

	// Authenication routes.
	mux.Handle("GET /signin", HandleSignIn(enableDebugAuth))
	mux.Handle("GET /github/signin", HandleOAuthSignIn(&githubConf))
	mux.Handle("GET /github/callback", HandleOAuthCallback(&githubConf, repo, FetchGithubUserID))
	mux.Handle("GET /google/signin", HandleOAuthSignIn(&googleConf))
	mux.Handle("GET /google/callback", HandleOAuthCallback(&googleConf, repo, FetchGoogleUserID))
	mux.Handle("POST /signout", HandleSignOutForm(repo))

	// Public blog routes.
	mux.Handle("GET /blogs", requireAccount(HandleBlogList(find)))
	mux.Handle("POST /blogs/create", requireAccount(HandleBlogCreateForm(repo, find, syncService)))
	mux.Handle("POST /blogs/{blogID}/follow", requireAccount(HandleBlogFollowForm(repo, find)))
	mux.Handle("POST /blogs/{blogID}/unfollow", requireAccount(HandleBlogUnfollowForm(repo, find)))

	// Public page routes.
	mux.Handle("GET /pages", requireAccount(HandlePageList(repo)))
	mux.Handle("POST /pages/create", requireAccount(HandlePageCreateForm(repo, pageFetcher)))
	mux.Handle("POST /pages/{pageID}/unfollow", requireAccount(HandlePageUnfollowForm(repo)))

	// Private (admin only) blog + post routes.
	mux.Handle("GET /blogs/{blogID}", requireAdmin(HandleBlogRead(repo)))
	mux.Handle("POST /blogs/{blogID}/delete", requireAdmin(HandleBlogDeleteForm(repo)))
	mux.Handle("GET /blogs/{blogID}/posts/{postID}", requireAdmin(HandlePostRead(repo)))
	mux.Handle("POST /blogs/{blogID}/posts/{postID}/delete", requireAdmin(HandlePostDeleteForm(repo)))

	// Private (admin only) account routes.
	mux.Handle("GET /accounts", requireAdmin(HandleAccountList(repo)))
	mux.Handle("POST /accounts/{accountID}/delete", requireAdmin(HandleAccountDeleteForm(repo)))

	// Debug endpoint for testing toasts.
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
		middleware.AddConfig(conf),
		middleware.CompressFiles(),
		middleware.PreventCSRF(),
		middleware.AddSecureHeaders(),
		middleware.LimitRequestBodySize(),
		middleware.Authenticate(repo),
	)

	return handler
}
