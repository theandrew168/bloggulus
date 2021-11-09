package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/theandrew168/bloggulus/internal/core"
)

type Application struct {
	blogStorage core.BlogStorage
	postStorage core.PostStorage
	logger      *log.Logger
}

func NewApplication(blogStorage core.BlogStorage, postStorage core.PostStorage, logger *log.Logger) *Application {
	app := Application{
		blogStorage: blogStorage,
		postStorage: postStorage,
		logger:      logger,
	}
	return &app
}

func (app *Application) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello API!"))
	})
	return r
}
