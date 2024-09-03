package page

import (
	_ "embed"

	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed signin.html
var SigninHTML string

type SigninData struct {
	layout.BaseData

	Username string
	Errors   map[string]string
}
