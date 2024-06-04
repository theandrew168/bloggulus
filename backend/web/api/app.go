package api

import (
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/api/admin"
	"github.com/theandrew168/bloggulus/backend/web/api/middleware"
	"github.com/theandrew168/bloggulus/backend/web/api/reader"
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
	adminApp := admin.NewApplication(app.store)
	readerApp := reader.NewApplication(app.store)

	mux := flow.New()
	mux.NotFound = http.HandlerFunc(util.NotFoundResponse)
	mux.MethodNotAllowed = http.HandlerFunc(util.MethodNotAllowedResponse)

	mux.Use(middleware.SecureHeaders)
	mux.Use(middleware.EnableCORS)

	mux.HandleFunc("/", app.handleIndexRapidoc(), "GET")
	mux.HandleFunc("/redoc", app.handleIndexRedoc(), "GET")
	mux.HandleFunc("/rapidoc", app.handleIndexRapidoc(), "GET")
	mux.HandleFunc("/stoplight", app.handleIndexStoplight(), "GET")
	mux.Handle("/admin/...", http.StripPrefix("/admin", adminApp.Router()))
	mux.Handle("/...", readerApp.Router())

	return mux
}
