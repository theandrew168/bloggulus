package api

import (
	"net/http"

	"github.com/alexedwards/flow"

	adminStorage "github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	readerStorage "github.com/theandrew168/bloggulus/backend/domain/reader/storage"
	"github.com/theandrew168/bloggulus/backend/web/api/admin"
	"github.com/theandrew168/bloggulus/backend/web/api/reader"
	"github.com/theandrew168/bloggulus/backend/web/api/util"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
)

type Application struct {
	adminStorage  adminStorage.Storage
	readerStorage readerStorage.Storage
}

func NewApplication(adminStorage adminStorage.Storage, readerStorage readerStorage.Storage) *Application {
	app := Application{
		adminStorage:  adminStorage,
		readerStorage: readerStorage,
	}
	return &app
}

func (app *Application) Router() http.Handler {
	adminApp := admin.NewApplication(app.adminStorage)
	readerApp := reader.NewApplication(app.readerStorage)

	mux := flow.New()
	mux.NotFound = http.HandlerFunc(util.NotFoundResponse)
	mux.MethodNotAllowed = http.HandlerFunc(util.MethodNotAllowedResponse)

	mux.Use(middleware.SecureHeaders)
	mux.Use(middleware.EnableCORS)

	mux.HandleFunc("/", app.handleIndex(), "GET")
	mux.Handle("/admin/...", http.StripPrefix("/admin", adminApp.Router()))
	mux.Handle("/...", readerApp.Router())

	return mux
}
