package reader

import (
	"net/http"

	"github.com/theandrew168/bloggulus/backend/storage"
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

func (app *Application) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /posts", app.handlePostList())

	return mux
}
