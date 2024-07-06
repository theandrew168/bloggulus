package api

import (
	"net/http"

	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
)

func Handler(store *storage.Storage, syncService *service.SyncService) http.Handler {
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

	mux.Handle("GET /articles", HandleArticleList(store))

	mux.Handle("POST /blogs", HandleBlogCreate(store, syncService))
	mux.Handle("GET /blogs", HandleBlogList(store))
	mux.Handle("GET /blogs/{blogID}", HandleBlogRead(store))
	mux.Handle("DELETE /blogs/{blogID}", HandleBlogDelete(store))

	mux.Handle("GET /blogs/{blogID}/posts", HandlePostList(store))
	mux.Handle("GET /blogs/{blogID}/posts/{postID}", HandlePostRead(store))
	mux.Handle("DELETE /blogs/{blogID}/posts/{postID}", HandlePostDelete(store))

	mux.Handle("POST /tags", HandleTagCreate(store))
	mux.Handle("GET /tags", HandleTagList(store))
	mux.Handle("DELETE /tags/{tagID}", HandleTagDelete(store))

	mux.Handle("POST /accounts", HandleAccountCreate(store))

	mux.Handle("POST /tokens", HandleTokenCreate(store))

	return middleware.Use(mux,
		middleware.SecureHeaders(),
		middleware.EnableCORS(),
		middleware.Authenticate(store),
	)
}
