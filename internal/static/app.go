package static

import (
	"embed"
	"io/fs"
	"net/http"
	"os"

	"github.com/klauspost/compress/gzhttp"
)

//go:embed static/img/logo.webp
var Favicon []byte

//go:embed static/etc/robots.txt
var Robots []byte

//go:embed static
var staticFS embed.FS

type Application struct {
	static fs.FS
}

func NewApplication() *Application {
	var static fs.FS
	if os.Getenv("DEBUG") != "" {
		// reload static files from filesystem if var DEBUG is set
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		static = os.DirFS("./internal/static/static/")
	} else {
		// else use the embedded static dir
		var err error
		static, err = fs.Sub(staticFS, "static")
		if err != nil {
			panic(err)
		}
	}

	app := Application{
		static: static,
	}

	return &app
}

func (app *Application) Router() http.Handler {
	// setup automatic compression handler for static files
	staticServer := http.FileServer(http.FS(app.static))
	gzipStaticServer := gzhttp.GzipHandler(staticServer)
	return gzipStaticServer
}
