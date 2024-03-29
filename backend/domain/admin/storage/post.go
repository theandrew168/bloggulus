package storage

import (
	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
)

type PostStorage interface {
	Create(post admin.Post) error
	Read(id uuid.UUID) (admin.Post, error)
	List(limit, offset int) ([]admin.Post, error)
	ListByBlog(blog admin.Blog, limit, offset int) ([]admin.Post, error)
}
