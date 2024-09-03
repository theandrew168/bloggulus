package page

import (
	_ "embed"
)

//go:embed signin.html
var SigninHTML string

type SigninData struct {
	Username string
	Errors   map[string]string
}
