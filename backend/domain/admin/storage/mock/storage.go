package mock

import "github.com/theandrew168/bloggulus/backend/domain/admin/storage"

// ensure Storage interface is satisfied
var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	mockBlogStorage *BlogStorage
	mockPostStorage *PostStorage
	mockTagStorage  *TagStorage
}

func NewStorage() *Storage {
	s := Storage{
		mockBlogStorage: NewBlogStorage(),
		mockPostStorage: NewPostStorage(),
		mockTagStorage:  NewTagStorage(),
	}
	return &s
}

func (s *Storage) Blog() storage.BlogStorage {
	return s.mockBlogStorage
}

func (s *Storage) Post() storage.PostStorage {
	return s.mockPostStorage
}

func (s *Storage) Tag() storage.TagStorage {
	return s.mockTagStorage
}

func (s *Storage) WithTransaction(operation func(store storage.Storage) error) error {
	return operation(s)
}
