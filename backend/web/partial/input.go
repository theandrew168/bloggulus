package partial

import (
	_ "embed"
)

//go:embed input.html
var InputHTML string

type InputData struct {
	Type        string
	Name        string
	Value       string
	Placeholder string
}
