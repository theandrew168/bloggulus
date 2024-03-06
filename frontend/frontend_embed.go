//go:build embed

package frontend

import "embed"

// embed with all:build to pickup _app dir (starts with an underscore)

//go:embed all:build
var frontend embed.FS

func init() {
	IsEmbedded = true
}
