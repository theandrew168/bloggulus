package web

import (
	"io/fs"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	"github.com/theandrew168/bloggulus/backend/command"
	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/jobqueue"
	webquery "github.com/theandrew168/bloggulus/backend/query/web"
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
	cmd *command.Command,
	qry *webquery.Query,
	queue *jobqueue.Queue,
) http.Handler {
	githubConf := oauth2.Config{
		Endpoint:     github.Endpoint,
		ClientID:     conf.GithubOAuth.ClientID,
		ClientSecret: conf.GithubOAuth.ClientSecret,
		RedirectURL:  conf.GithubOAuth.RedirectURI,
		Scopes:       []string{},
	}
	googleConf := oauth2.Config{
		Endpoint:     google.Endpoint,
		ClientID:     conf.GoogleOAuth.ClientID,
		ClientSecret: conf.GoogleOAuth.ClientSecret,
		RedirectURL:  conf.GoogleOAuth.RedirectURI,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
	}

	rootMux := http.NewServeMux()
	webMux := http.NewServeMux()

	requireAccount := middleware.RequireAccount()
	requireAdmin := middleware.Chain(requireAccount, middleware.RequireAdmin())

	// Host prometheus metrics on "/metrics".
	webMux.Handle("GET /metrics", promhttp.Handler())

	// Basic health check endpoint for verifying that the app is running.
	webMux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	})

	// Compress and serve the embedded public (static) files.
	publicFiles, _ := fs.Sub(public, "public")
	publicFilesHandler := http.FileServer(http.FS(publicFiles))

	// Serve public (static) files from the embedded FS.
	webMux.Handle("GET /favicon.ico", publicFilesHandler)
	webMux.Handle("GET /robots.txt", publicFilesHandler)
	webMux.Handle("GET /css/", publicFilesHandler)
	webMux.Handle("GET /img/", publicFilesHandler)
	webMux.Handle("GET /js/", publicFilesHandler)

	// The main application routes start here.
	webMux.Handle("GET /{$}", HandleIndexPage(qry))

	// Check if the debug auth method should be enabled.
	if conf.EnableDebugAuth {
		webMux.Handle("POST /signin/debug", HandleDebugSignIn(conf.SecretKey, cmd))
	}

	// Authenication routes.
	webMux.Handle("GET /signin", HandleSignIn(conf.EnableDebugAuth))
	webMux.Handle("GET /signin/github", HandleOAuthSignIn(&githubConf))
	webMux.Handle("GET /signin/github/callback", HandleOAuthCallback(conf.SecretKey, cmd, &githubConf, FetchGithubUserID))
	webMux.Handle("GET /signin/google", HandleOAuthSignIn(&googleConf))
	webMux.Handle("GET /signin/google/callback", HandleOAuthCallback(conf.SecretKey, cmd, &googleConf, FetchGoogleUserID))
	webMux.Handle("POST /signout", HandleSignOutForm(cmd))

	// Documentation routes.
	webMux.Handle("GET /docs/privacy", HandlePrivacy())

	// Public blog routes.
	webMux.Handle("GET /blogs", requireAccount(HandleBlogList(qry)))
	webMux.Handle("POST /blogs/create", requireAccount(HandleBlogCreateForm(cmd, qry, queue)))
	webMux.Handle("POST /blogs/{blogID}/follow", requireAccount(HandleBlogFollowForm(cmd, qry)))
	webMux.Handle("POST /blogs/{blogID}/unfollow", requireAccount(HandleBlogUnfollowForm(cmd, qry)))

	// Private (admin only) blog + post routes.
	webMux.Handle("GET /blogs/{blogID}", requireAdmin(HandleBlogRead(qry)))
	webMux.Handle("POST /blogs/{blogID}/delete", requireAdmin(HandleBlogDeleteForm(cmd)))
	webMux.Handle("POST /blogs/{blogID}/hide", requireAdmin(HandleBlogHideForm(cmd)))
	webMux.Handle("POST /blogs/{blogID}/show", requireAdmin(HandleBlogShowForm(cmd)))
	webMux.Handle("GET /blogs/{blogID}/posts/{postID}", requireAdmin(HandlePostRead(qry)))
	webMux.Handle("POST /blogs/{blogID}/posts/{postID}/delete", requireAdmin(HandlePostDeleteForm(cmd)))

	// Private (admin only) account routes.
	webMux.Handle("GET /accounts", requireAdmin(HandleAccountList(qry)))
	webMux.Handle("POST /accounts/{accountID}/delete", requireAdmin(HandleAccountDeleteForm(cmd)))

	// Private (admin only) maintenance routes.
	webMux.Handle("GET /jobs", requireAdmin(HandleJobsList(queue)))
	webMux.Handle("POST /jobs/{jobID}/delete", requireAdmin(HandleJobDeleteForm(queue)))

	// Debug endpoint for testing toasts.
	webMux.HandleFunc("GET /toast", func(w http.ResponseWriter, r *http.Request) {
		cookie := util.NewSessionCookie(util.ToastCookieName, "Toasts are awesome!")
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Apply api-specific middleware to the api mux.
	apiHandler := middleware.Use(api.Handler(qry),
		// The trace middleware needs to be as close to the ServeMux as possible
		// in order to capture r.Pattern before any other middleware calls r.WithContext.
		// Putting it last here means that it will run directly before and directly after
		// the ServeMux itself.
		middleware.Trace("/api/v1"),
	)
	rootMux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiHandler))

	// Requests that don't match any of the above handlers get a 404.
	webMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		util.NotFoundResponse(w, r)
	})

	// Apply web-specific middleware to the web mux.
	webHandler := middleware.Use(webMux,
		middleware.PreventCSRF(),
		middleware.Authenticate(qry),
		// The trace middleware needs to be as close to the ServeMux as possible
		// in order to capture r.Pattern before any other middleware calls r.WithContext.
		// Putting it last here means that it will run directly before and directly after
		// the ServeMux itself.
		middleware.Trace(""),
	)
	rootMux.Handle("/", webHandler)

	// Apply global middleware to the root mux.
	rootHandler := middleware.Use(rootMux,
		middleware.RecoverPanic(),
		middleware.CompressFiles(),
		middleware.AddSecureHeaders(),
		middleware.LimitRequestBodySize(),
	)

	return rootHandler
}
