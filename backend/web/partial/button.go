package partial

import (
	_ "embed"
)

//go:embed button.html
var ButtonHTML string

type ButtonData struct {
	Contents string
}
