package api

import (
	"github.com/theandrew168/bloggulus/internal/core"
)

type Application struct {
	Blog    core.BlogStorage
	Post    core.PostStorage
}
