package storage

import (
	"github.com/theandrew168/bloggulus/backend/domain/reader"
)

type PostStorage interface {
	List(limit, offset int) ([]*reader.Post, error)
	Search(query string, limit, offset int) ([]*reader.Post, error)
}
