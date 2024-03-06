//go:build !embed

package frontend

import "embed"

var frontend embed.FS

func init() {
	IsEmbedded = false
}
