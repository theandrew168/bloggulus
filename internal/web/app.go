package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

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
	var templates fs.FS
	if strings.HasPrefix(cfg.Env, "dev") {
		// reload templates from filesystem if Config.Env starts with "dev"
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		templates = os.DirFS("./internal/web/templates/")
	} else {
		// else use the embedded templates dir
		templates, _ = fs.Sub(templatesFS, "templates")
	}

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
	return r
}
