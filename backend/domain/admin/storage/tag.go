package storage

import (
	"github.com/theandrew168/bloggulus/backend/domain/admin"
)

type TagStorage interface {
	Create(tag admin.Tag) error
	List(limit, offset int) ([]admin.Tag, error)
	Delete(tag admin.Tag) error
}
