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

// ensure TagStorage interface is satisfied
var _ TagStorage = (*MockTagStorage)(nil)

type MockTagStorage struct {
	data map[uuid.UUID]*admin.Tag
}

func NewMockTagStorage() *MockTagStorage {
	s := MockTagStorage{
		data: make(map[uuid.UUID]*admin.Tag),
	}
	return &s
}

func (s *MockTagStorage) Create(tag *admin.Tag) error {
	_, ok := s.data[tag.ID()]
	if ok {
		return ErrConflict
	}

	s.data[tag.ID()] = tag
	return nil
}

func (s *MockTagStorage) Read(id uuid.UUID) (*admin.Tag, error) {
	tag, ok := s.data[id]
	if !ok {
		return nil, ErrNotFound
	}

	return tag, nil
}

func (s *MockTagStorage) List(limit, offset int) ([]*admin.Tag, error) {
	var tags []*admin.Tag
	for _, tag := range s.data {
		tags = append(tags, tag)
	}

	start := offset
	end := offset + limit
	if start >= len(tags) || end >= len(tags) {
		return nil, nil
	}

	return tags[start:end], nil
}

func (s *MockTagStorage) Delete(tag *admin.Tag) error {
	_, ok := s.data[tag.ID()]
	if !ok {
		return ErrNotFound
	}

	delete(s.data, tag.ID())
	return nil
}
