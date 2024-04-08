package storage

import (
	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
)

type TagStorage interface {
	Create(tag *admin.Tag) error
	Read(id uuid.UUID) (*admin.Tag, error)
	List(limit, offset int) ([]*admin.Tag, error)
	Delete(tag *admin.Tag) error
}
