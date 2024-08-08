package api

import (
	"net/http"

	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
)

func Handler(store *storage.Storage, syncService *service.SyncService) http.Handler {
	mux := http.NewServeMux()

	accountRequired := middleware.AccountRequired()
	adminRequired := middleware.Chain(accountRequired, middleware.AdminRequired())

	mux.Handle("GET /{$}", HandleIndexRapidoc())
	mux.Handle("GET /redoc", HandleIndexRedoc())
	mux.Handle("GET /rapidoc", HandleIndexRapidoc())
	mux.Handle("GET /stoplight", HandleIndexStoplight())

	mux.Handle("GET /articles", HandleArticleList(store))

	mux.Handle("POST /blogs", adminRequired(HandleBlogCreate(store, syncService)))
	mux.Handle("GET /blogs", accountRequired(HandleBlogList(store)))
	mux.Handle("GET /blogs/{blogID}", accountRequired(HandleBlogRead(store)))
	mux.Handle("DELETE /blogs/{blogID}", adminRequired(HandleBlogDelete(store)))

	// Follow a blog: HandleBlogFollow
	// POST /blogs/{blogID}/follow -> 204
	mux.Handle("POST /blogs/{blogID}/follow", accountRequired(HandleBlogFollow(store)))

	// Unfollow a blog: HandleBlogUnfollow
	// POST /blogs/{blogID}/unfollow -> 204
	mux.Handle("POST /blogs/{blogID}/unfollow", accountRequired(HandleBlogUnfollow(store)))

	// See what blogs the auth'd account follows (GH style): HandleBlogListFollowing
	// GET /blogs/following -> []Blog
	mux.Handle("GET /blogs/following", accountRequired(HandleBlogListFollowing(store)))

	mux.Handle("GET /blogs/{blogID}/posts", accountRequired(HandlePostList(store)))
	mux.Handle("GET /blogs/{blogID}/posts/{postID}", accountRequired(HandlePostRead(store)))
	mux.Handle("DELETE /blogs/{blogID}/posts/{postID}", adminRequired(HandlePostDelete(store)))

	mux.Handle("POST /tags", adminRequired(HandleTagCreate(store)))
	mux.Handle("GET /tags", accountRequired(HandleTagList(store)))
	mux.Handle("DELETE /tags/{tagID}", adminRequired(HandleTagDelete(store)))

	mux.Handle("POST /accounts", HandleAccountCreate(store))

	mux.Handle("POST /tokens", HandleTokenCreate(store))

	return middleware.Use(mux,
		middleware.SecureHeaders(),
		middleware.EnableCORS(),
		middleware.Authenticate(store),
	)
}
