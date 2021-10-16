package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/theandrew168/bloggulus/internal/core"
)

type Application struct {
	Blog core.BlogStorage
	Post core.PostStorage
}

func (app *Application) Router() http.Handler {
	router := httprouter.New()
	router.HandlerFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello API!"))
	})
	return router
}
