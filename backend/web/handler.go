package web

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	"github.com/theandrew168/bloggulus/backend/command"
	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/job"
	"github.com/theandrew168/bloggulus/backend/query"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/web/api"
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
	cmd *command.Command,
	qry *query.Query,
	syncService *job.SyncService,
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
	mux.Handle("/docs/", publicFilesHandler)
	mux.Handle("/css/", publicFilesHandler)
	mux.Handle("/img/", publicFilesHandler)
	mux.Handle("/js/", publicFilesHandler)

	// The main application routes start here.
	mux.Handle("GET /{$}", HandleIndexPage(qry))

	apiHandler := api.Handler(public, cmd, qry)
	mux.Handle("GET /api/v1/", http.StripPrefix("/api/v1", apiHandler))

	// Check if the debug auth method should be enabled.
	enableDebugAuth := os.Getenv("ENABLE_DEBUG_AUTH") != ""
	if enableDebugAuth {
		mux.Handle("POST /debug/signin", HandleDebugSignIn(conf.SecretKey, cmd))
	}

	// Authenication routes.
	mux.Handle("GET /signin", HandleSignIn(enableDebugAuth))
	mux.Handle("GET /github/signin", HandleOAuthSignIn(&githubConf))
	mux.Handle("GET /github/callback", HandleOAuthCallback(conf.SecretKey, cmd, &githubConf, FetchGithubUserID))
	mux.Handle("GET /google/signin", HandleOAuthSignIn(&googleConf))
	mux.Handle("GET /google/callback", HandleOAuthCallback(conf.SecretKey, cmd, &googleConf, FetchGoogleUserID))
	mux.Handle("POST /signout", HandleSignOutForm(cmd))

	// Public blog routes.
	mux.Handle("GET /blogs", requireAccount(HandleBlogList(qry)))
	mux.Handle("POST /blogs/create", requireAccount(HandleBlogCreateForm(repo, syncService)))
	mux.Handle("POST /blogs/{blogID}/follow", requireAccount(HandleBlogFollowForm(repo)))
	mux.Handle("POST /blogs/{blogID}/unfollow", requireAccount(HandleBlogUnfollowForm(repo)))

	// Private (admin only) blog + post routes.
	mux.Handle("GET /blogs/{blogID}", requireAdmin(HandleBlogRead(repo)))
	mux.Handle("POST /blogs/{blogID}/delete", requireAdmin(HandleBlogDeleteForm(cmd)))
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
