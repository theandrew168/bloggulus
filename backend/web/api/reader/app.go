package reader

import (
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/bloggulus/backend/domain/reader/storage"
	"github.com/theandrew168/bloggulus/backend/web/api/util"
)

type Application struct {
	storage storage.Storage
}

func NewApplication(storage storage.Storage) *Application {
	app := Application{
		storage: storage,
	}
	return &app
}

func (app *Application) Router() http.Handler {
	mux := flow.New()
	mux.NotFound = http.HandlerFunc(util.NotFoundResponse)
	mux.MethodNotAllowed = http.HandlerFunc(util.MethodNotAllowedResponse)

	mux.HandleFunc("/posts", app.handlePostList(), "GET")

	return mux
}
