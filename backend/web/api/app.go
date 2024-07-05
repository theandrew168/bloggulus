package api

import (
	"net/http"

	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
)

type Application struct {
	store *storage.Storage

	syncService *service.SyncService
}

func NewApplication(store *storage.Storage, syncService *service.SyncService) *Application {
	app := Application{
		store: store,

		syncService: syncService,
	}
	return &app
}

func (app *Application) Handler() http.Handler {
	mux := http.NewServeMux()

	// accountRequired := middleware.AccountRequired()
	// adminRequired := middleware.AdminRequired()
	// protected := middleware.Chain(accountRequired, adminRequired)

	mux.Handle("GET /{$}", app.handleIndexRapidoc())
	mux.Handle("GET /redoc", app.handleIndexRedoc())
	mux.Handle("GET /rapidoc", app.handleIndexRapidoc())
	mux.Handle("GET /stoplight", app.handleIndexStoplight())

	mux.Handle("GET /articles", app.handleArticleList())

	mux.Handle("POST /blogs", app.handleBlogCreate())
	mux.Handle("GET /blogs", app.handleBlogList())
	mux.Handle("GET /blogs/{blogID}", app.handleBlogRead())
	mux.Handle("DELETE /blogs/{blogID}", app.handleBlogDelete())

	mux.Handle("GET /blogs/{blogID}/posts", app.handlePostList())
	mux.Handle("GET /blogs/{blogID}/posts/{postID}", app.handlePostRead())
	mux.Handle("DELETE /blogs/{blogID}/posts/{postID}", app.handlePostDelete())

	mux.Handle("POST /tags", app.handleTagCreate())
	mux.Handle("GET /tags", app.handleTagList())
	mux.Handle("DELETE /tags/{tagID}", app.handleTagDelete())

	mux.Handle("POST /accounts", app.handleAccountCreate())

	mux.Handle("POST /tokens", app.handleTokenCreate())

	return middleware.Use(mux,
		middleware.SecureHeaders(),
		middleware.EnableCORS(),
		middleware.Authenticate(app.store),
	)
}
