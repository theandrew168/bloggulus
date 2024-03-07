package api

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
)

//go:embed templates
var templatesFS embed.FS

type Application struct {
	templates fs.FS

	logger  *log.Logger
	storage *storage.Storage
}

func NewApplication(logger *log.Logger, storage *storage.Storage) *Application {
	templates, err := fs.Sub(templatesFS, "templates")
	if err != nil {
		panic(err)
	}

	app := Application{
		templates: templates,

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

	mux.HandleFunc("/", app.HandleIndex, "GET")
	mux.HandleFunc("/blogs", app.HandleReadBlogs, "GET")
	mux.HandleFunc("/blogs/:id", app.HandleReadBlog, "GET")
	mux.HandleFunc("/posts", app.HandleReadPosts, "GET")
	mux.HandleFunc("/posts/:id", app.HandleReadPost, "GET")

	return mux
}
