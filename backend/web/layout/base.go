package layout

import (
	_ "embed"

	"github.com/theandrew168/bloggulus/backend/model"
)

//go:embed base.html
var BaseHTML string

type BaseData struct {
	Account   *model.Account
	CSRFToken string
	Toast     string

	GoatCounterCode string
}
