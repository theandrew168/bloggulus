package page

import (
	_ "embed"

	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed register.html
var RegisterHTML string

type RegisterData struct {
	layout.BaseData

	Username string
	Errors   map[string]string
}
