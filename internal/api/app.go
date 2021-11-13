package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/theandrew168/bloggulus/internal/core"
)

type Application struct {
	storage core.Storage
	logger  *log.Logger
}

func NewApplication(storage core.Storage, logger *log.Logger) *Application {
	app := Application{
		storage: storage,
		logger:  logger,
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
