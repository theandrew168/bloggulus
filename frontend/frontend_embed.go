//go:build embed

package frontend

import (
	"embed"
	"io/fs"
)

//go:embed dist
var frontend embed.FS

func init() {
	var err error
	Frontend, err = fs.Sub(frontend, "dist")
	if err != nil {
		panic(err)
	}
}
