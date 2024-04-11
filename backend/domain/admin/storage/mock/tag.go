package mock

import (
	"sync"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
)

// ensure TagStorage interface is satisfied
var _ storage.TagStorage = (*TagStorage)(nil)

type TagStorage struct {
	mu   sync.RWMutex
	data map[uuid.UUID]*admin.Tag
}

func NewTagStorage() *TagStorage {
	s := TagStorage{
		data: make(map[uuid.UUID]*admin.Tag),
	}
	return &s
}

func (s *TagStorage) Create(tag *admin.Tag) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[tag.ID()]
	if ok {
		return storage.ErrConflict
	}

	s.data[tag.ID()] = tag
	return nil
}

func (s *TagStorage) Read(id uuid.UUID) (*admin.Tag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tag, ok := s.data[id]
	if !ok {
		return nil, storage.ErrNotFound
	}

	return tag, nil
}

func (s *TagStorage) List(limit, offset int) ([]*admin.Tag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tags []*admin.Tag
	for _, tag := range s.data {
		tags = append(tags, tag)
	}

	start := offset
	end := min(offset+limit, len(tags))
	if start >= len(tags) {
		return nil, nil
	}

	return tags[start:end], nil
}

func (s *TagStorage) Delete(tag *admin.Tag) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[tag.ID()]
	if !ok {
		return storage.ErrNotFound
	}

	delete(s.data, tag.ID())
	return nil
}
