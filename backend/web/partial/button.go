package partial

import (
	_ "embed"
)

//go:embed button.html
var ButtonHTML string

type ButtonData struct {
	Name     string
	Value    string
	Contents string
}
