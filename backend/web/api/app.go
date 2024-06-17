package api

import (
	"net/http"

	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/api/admin"
	"github.com/theandrew168/bloggulus/backend/web/api/reader"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
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
	adminApp := admin.NewApplication(app.store)
	readerApp := reader.NewApplication(app.store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", app.handleIndexRapidoc())
	mux.HandleFunc("GET /redoc", app.handleIndexRedoc())
	mux.HandleFunc("GET /rapidoc", app.handleIndexRapidoc())
	mux.HandleFunc("GET /stoplight", app.handleIndexStoplight())
	mux.Handle("/admin/", http.StripPrefix("/admin", adminApp.Router()))
	mux.Handle("/", readerApp.Router())

	return middleware.Adapt(mux, middleware.SecureHeaders(), middleware.EnableCORS())
}
