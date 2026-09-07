package layout

import (
	_ "embed"

	"github.com/theandrew168/bloggulus/backend/query"
)

//go:embed base.html
var BaseHTML string

type BaseData struct {
	Account     query.Account
	CSRFToken   string
	Toast       string
	Stylesheets []string
}
