package admin

import (
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/api/util"
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

func (app *Application) Router() http.Handler {
	mux := flow.New()
	mux.NotFound = http.HandlerFunc(util.NotFoundResponse)
	mux.MethodNotAllowed = http.HandlerFunc(util.MethodNotAllowedResponse)

	mux.HandleFunc("/blogs", app.handleBlogList(), "GET")
	mux.HandleFunc("/blogs/:id", app.handleBlogRead(), "GET")
	mux.HandleFunc("/posts", app.handlePostList(), "GET")
	mux.HandleFunc("/posts/:id", app.handlePostRead(), "GET")
	mux.HandleFunc("/tags", app.handleTagList(), "GET")

	return mux
}
