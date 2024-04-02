package storage

import (
	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
)

type BlogStorage interface {
	Create(blog *admin.Blog) error
	Read(id uuid.UUID) (*admin.Blog, error)
	ReadByFeedURL(feedURL string) (*admin.Blog, error)
	List(limit, offset int) ([]*admin.Blog, error)
	Update(blog *admin.Blog) error
	Delete(blog *admin.Blog) error
}
