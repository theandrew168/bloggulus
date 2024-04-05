package storage

type Storage interface {
	Blog() BlogStorage
	Post() PostStorage
	Tag() TagStorage

	WithTransaction(operation func(store Storage) error) error
}

// ensure Storage interface is satisfied
var _ Storage = (*MockStorage)(nil)

type MockStorage struct {
	mockBlogStorage *MockBlogStorage
	mockPostStorage *MockPostStorage
	mockTagStorage  *MockTagStorage
}

func NewMockStorage() *MockStorage {
	s := MockStorage{
		mockBlogStorage: NewMockBlogStorage(),
		mockPostStorage: NewMockPostStorage(),
		mockTagStorage:  NewMockTagStorage(),
	}
	return &s
}

func (s *MockStorage) Blog() BlogStorage {
	return s.mockBlogStorage
}

func (s *MockStorage) Post() PostStorage {
	return s.mockPostStorage
}

func (s *MockStorage) Tag() TagStorage {
	return s.mockTagStorage
}

func (s *MockStorage) WithTransaction(operation func(store Storage) error) error {
	return operation(s)
}
