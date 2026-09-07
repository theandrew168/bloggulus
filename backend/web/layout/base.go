package layout

import (
	_ "embed"

	webquery "github.com/theandrew168/bloggulus/backend/query/web"
)

//go:embed base.html
var BaseHTML string

type BaseData struct {
	Account     webquery.Account
	CSRFToken   string
	Toast       string
	Stylesheets []string
}
