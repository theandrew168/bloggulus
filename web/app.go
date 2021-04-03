package web

import (
	"github.com/theandrew168/bloggulus/storage"
)

type Application struct {
	Account     storage.Account
	Blog        storage.Blog
	Post        storage.Post
	Session     storage.Session
}
