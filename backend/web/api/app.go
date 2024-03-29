package api

import (
	"log"
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
)

type Application struct {
	logger  *log.Logger
	storage storage.Storage
}

func NewApplication(logger *log.Logger, storage storage.Storage) *Application {
	app := Application{
		logger:  logger,
		storage: storage,
	}
	return &app
}

func (app *Application) Router() http.Handler {
	mux := flow.New()
	mux.NotFound = http.HandlerFunc(app.notFoundResponse)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	mux.Use(middleware.SecureHeaders)
	mux.Use(middleware.EnableCORS)

	mux.HandleFunc("/", app.handleIndex(), "GET")
	mux.HandleFunc("/blogs", app.handleBlogList(), "GET")
	mux.HandleFunc("/blogs/:id", app.handleBlogRead(), "GET")
	mux.HandleFunc("/posts", app.handlePostList(), "GET")
	mux.HandleFunc("/posts/:id", app.handlePostRead(), "GET")
	mux.HandleFunc("/tags", app.handleTagList(), "GET")

	return mux
}
