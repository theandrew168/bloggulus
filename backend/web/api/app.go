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

	// TODO: Three options for testing the handlers indepedently of auth middleware:
	// 1. Add auth to the tests for each route (bad idea / couples too many ideas together)
	// 2. Don't test against the top-level app.Handler(), test specific handlers (need to make them public)
	// 3. De-couple handlers from the joint App (closure in the required store / services)
	// Let's go with number 3!

	// accountRequired := middleware.AccountRequired()
	// adminRequired := middleware.AdminRequired()
	// protected := middleware.Chain(accountRequired, adminRequired)

	mux.Handle("GET /{$}", HandleIndexRapidoc())
	mux.Handle("GET /redoc", HandleIndexRedoc())
	mux.Handle("GET /rapidoc", HandleIndexRapidoc())
	mux.Handle("GET /stoplight", HandleIndexStoplight())

	mux.Handle("GET /articles", HandleArticleList(app.store))

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
