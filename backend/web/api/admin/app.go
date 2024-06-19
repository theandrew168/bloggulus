package admin

import (
	"net/http"

	"github.com/theandrew168/bloggulus/backend/storage"
)

type Application struct {
	store *storage.Storage
}

func NewApplication(store *storage.Storage) *Application {
	app := Application{
		store: store,
	}
	return &app
}

func (app *Application) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /blogs", app.handleBlogList())
	mux.HandleFunc("GET /blogs/{id}", app.handleBlogRead())
	mux.HandleFunc("GET /posts", app.handlePostList())
	mux.HandleFunc("GET /posts/{id}", app.handlePostRead())
	mux.HandleFunc("GET /tags", app.handleTagList())
	mux.HandleFunc("POST /accounts", app.handleAccountCreate())
	mux.HandleFunc("POST /tokens", app.handleTokenCreate())

	return mux
}
