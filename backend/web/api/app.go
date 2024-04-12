package api

import (
	"net/http"

	"github.com/alexedwards/flow"

	adminStorage "github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/web/api/admin"
	"github.com/theandrew168/bloggulus/backend/web/api/util"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
)

type Application struct {
	adminStorage adminStorage.Storage
}

func NewApplication(adminStorage adminStorage.Storage) *Application {
	app := Application{
		adminStorage: adminStorage,
	}
	return &app
}

func (app *Application) Router() http.Handler {
	adminApp := admin.NewApplication(app.adminStorage)

	mux := flow.New()
	mux.NotFound = http.HandlerFunc(util.NotFoundResponse)
	mux.MethodNotAllowed = http.HandlerFunc(util.MethodNotAllowedResponse)

	mux.Use(middleware.SecureHeaders)
	mux.Use(middleware.EnableCORS)

	mux.HandleFunc("/", app.handleIndex(), "GET")
	mux.Handle("/admin/...", http.StripPrefix("/admin", adminApp.Router()))

	return mux
}
