package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/bloggulus/internal/middleware"
	"github.com/theandrew168/bloggulus/internal/storage"
)

var (
	pageSize = 15
)

//go:embed template
var templateFS embed.FS

type Application struct {
	templates fs.FS

	logger  *log.Logger
	storage *storage.Storage
}

func NewApplication(logger *log.Logger, storage *storage.Storage) *Application {
	var templates fs.FS
	if os.Getenv("DEBUG") != "" {
		// reload templates from filesystem if var DEBUG is set
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		templates = os.DirFS("./internal/web/template/")
	} else {
		// else use the embedded template dir
		templates, _ = fs.Sub(templateFS, "template")
	}

	app := Application{
		templates: templates,

		logger:  logger,
		storage: storage,
	}
	return &app
}

func (app *Application) Router() http.Handler {
	mux := flow.New()
	mux.NotFound = http.HandlerFunc(app.notFoundResponse)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	mux.Use(middleware.RecoverPanic)
	mux.Use(middleware.SecureHeaders)

	mux.HandleFunc("/", app.HandleIndex, "GET")

	return mux
}
