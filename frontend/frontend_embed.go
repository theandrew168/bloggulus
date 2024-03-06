//go:build embed

package frontend

import "embed"

//go:embed all:build
var frontend embed.FS

func init() {
	IsEmbedded = true
}
