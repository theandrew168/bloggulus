package api

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/theandrew168/bloggulus/internal/config"
	"github.com/theandrew168/bloggulus/internal/core"
)

//go:embed templates
var templatesFS embed.FS

type Application struct {
	templates fs.FS
	storage   core.Storage
	logger    *log.Logger
	cfg       config.Config
}

func NewApplication(storage core.Storage, logger *log.Logger, cfg config.Config) *Application {
	templates, _ := fs.Sub(templatesFS, "templates")

	app := Application{
		templates: templates,
		storage:   storage,
		logger:    logger,
		cfg:       cfg,
	}
	return &app
}

func (app *Application) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", app.HandleIndex)
	r.Get("/blog", app.HandleBlog)
	return app.enableCORS(r)
}
