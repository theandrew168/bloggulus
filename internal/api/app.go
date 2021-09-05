package api

import (
	"github.com/theandrew168/bloggulus/internal/core"
)

type Application struct {
	Account     core.AccountStorage
	AccountBlog core.AccountBlogStorage
	Blog        core.BlogStorage
	Post        core.PostStorage
	Session     core.SessionStorage
}
