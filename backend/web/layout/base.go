package layout

import (
	_ "embed"
)

//go:embed base.html
var BaseHTML string

type BaseData struct{}
