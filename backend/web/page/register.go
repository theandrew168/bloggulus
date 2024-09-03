package page

import (
	_ "embed"
)

//go:embed register.html
var RegisterHTML string

type RegisterData struct {
	Username string
	Errors   map[string]string
}
