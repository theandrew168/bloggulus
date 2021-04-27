package app

import (
	"github.com/theandrew168/bloggulus/storage"
)

type Application struct {
	Account     storage.Account
	Blog        storage.Blog
	AccountBlog storage.AccountBlog
	Post        storage.Post
	Session     storage.Session
}
