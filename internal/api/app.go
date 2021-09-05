package api

import (
	"github.com/theandrew168/bloggulus/internal/model"
)

type Application struct {
	Account     model.AccountStorage
	AccountBlog model.AccountBlogStorage
	Blog        model.BlogStorage
	Post        model.PostStorage
	Session     model.SessionStorage
}
