//go:build embed

package frontend

import (
	"embed"
	"io/fs"
)

// embed with all:build to pickup _app dir (starts with an underscore)

//go:embed all:build
var frontend embed.FS

func init() {
	var err error
	Frontend, err = fs.Sub(frontend, "build")
	if err != nil {
		panic(err)
	}
}
