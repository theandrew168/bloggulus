package web

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/theandrew168/bloggulus/internal/core"
)

type Application struct {
	templates   fs.FS
	blogStorage core.BlogStorage
	postStorage core.PostStorage
	logger      *log.Logger
}

func NewApplication(templates fs.FS, blogStorage core.BlogStorage, postStorage core.PostStorage, logger *log.Logger) *Application {
	app := Application{
		templates:   templates,
		blogStorage: blogStorage,
		postStorage: postStorage,
		logger:      logger,
	}
	return &app
}

func (app *Application) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", app.HandleIndex)
	return r
}
