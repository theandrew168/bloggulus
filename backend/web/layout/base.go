package layout

import (
	_ "embed"

	"github.com/theandrew168/bloggulus/backend/query"
)

//go:embed base.html
var BaseHTML string

// TODO: Explore supporting extra style sheets.

type BaseData struct {
	Account   query.Account
	CSRFToken string
	Toast     string
}
