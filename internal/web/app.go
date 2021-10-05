package web

import (
	"io/fs"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/bloggulus/internal/core"
)

type Application struct {
	StaticFS    fs.FS
	TemplatesFS fs.FS

	Blog    core.BlogStorage
	Post    core.PostStorage
}

func (app *Application) Router() http.Handler {
	router := httprouter.New()
	router.HandlerFunc("GET", "/", app.HandleIndex)
	router.Handler("GET", "/metrics", promhttp.Handler())
	router.ServeFiles("/static/*filepath", http.FS(app.StaticFS))
	return router
}
