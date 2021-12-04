package api

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/theandrew168/bloggulus/internal/core"
)

var (
	queryTimeout = 3 * time.Second
)

//go:embed templates
var templatesFS embed.FS

type Application struct {
	templates fs.FS
	storage   core.Storage
	logger    *log.Logger
}

func NewApplication(storage core.Storage, logger *log.Logger) *Application {
	templates, _ := fs.Sub(templatesFS, "templates")

	app := Application{
		templates: templates,
		storage:   storage,
		logger:    logger,
	}
	return &app
}

func (app *Application) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{}))
	r.Use(middleware.Recoverer)

	r.NotFound(app.notFoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	r.Get("/", app.HandleIndex)
	r.Get("/blog", app.HandleReadBlogs)
	r.Get("/blog/{id}", app.HandleReadBlog)
	r.Get("/post", app.HandleReadPosts)
	r.Get("/post/{id}", app.HandleReadPost)

	return r
}
