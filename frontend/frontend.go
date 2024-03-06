package frontend

import "embed"

var Frontend embed.FS

var IsEmbedded bool

// conditional embedding based on:
// https://github.com/golang/go/issues/44484#issuecomment-948137497
func init() {
	Frontend = frontend
}
