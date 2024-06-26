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

	mux.HandleFunc("GET /{$}", app.handleIndexRapidoc())
	mux.HandleFunc("GET /redoc", app.handleIndexRedoc())
	mux.HandleFunc("GET /rapidoc", app.handleIndexRapidoc())
	mux.HandleFunc("GET /stoplight", app.handleIndexStoplight())

	mux.HandleFunc("GET /articles", app.handleArticleList())

	mux.HandleFunc("POST /blogs", app.handleBlogCreate())
	mux.HandleFunc("GET /blogs", app.handleBlogList())
	mux.HandleFunc("GET /blogs/{blogID}", app.handleBlogRead())
	mux.HandleFunc("DELETE /blogs/{blogID}", app.handleBlogDelete())

	mux.HandleFunc("GET /blogs/{blogID}/posts", app.handlePostList())
	mux.HandleFunc("GET /blogs/{blogID}/posts/{postID}", app.handlePostRead())
	mux.HandleFunc("DELETE /blogs/{blogID}/posts/{postID}", app.handlePostDelete())

	mux.HandleFunc("POST /tags", app.handleTagCreate())
	mux.HandleFunc("GET /tags", app.handleTagList())
	mux.HandleFunc("DELETE /tags/{tagID}", app.handleTagDelete())

	mux.HandleFunc("POST /accounts", app.handleAccountCreate())

	mux.HandleFunc("POST /tokens", app.handleTokenCreate())

	return middleware.Use(mux,
		middleware.SecureHeaders(),
		middleware.EnableCORS(),
		middleware.Authenticate(app.store),
	)
}
