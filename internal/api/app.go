package api

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/theandrew168/bloggulus/internal/storage"
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
	if strings.HasPrefix(os.Getenv("ENV"), "dev") {
		// reload templates from filesystem if var ENV starts with "dev"
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		templates = os.DirFS("./internal/api/template/")
	} else {
		// else use the embedded template FS
		var err error
		templates, err = fs.Sub(templateFS, "template")
		if err != nil {
			panic(err)
		}
	}

	app := Application{
		templates: templates,

		logger:  logger,
		storage: storage,
	}
	return &app
}

func (app *Application) Router() http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{}))
	mux.Use(middleware.Recoverer)

	mux.NotFound(app.notFoundResponse)
	mux.MethodNotAllowed(app.methodNotAllowedResponse)

	mux.Get("/", app.HandleIndex)
	mux.Get("/blog", app.HandleReadBlogs)
	mux.Get("/blog/{id}", app.HandleReadBlog)
	mux.Get("/post", app.HandleReadPosts)
	mux.Get("/post/{id}", app.HandleReadPost)

	return mux
}
