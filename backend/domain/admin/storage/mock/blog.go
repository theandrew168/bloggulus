package mock

import (
	"sync"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
)

// ensure BlogStorage interface is satisfied
var _ storage.BlogStorage = (*BlogStorage)(nil)

type BlogStorage struct {
	mu   sync.RWMutex
	data map[uuid.UUID]*admin.Blog
}

func NewBlogStorage() *BlogStorage {
	s := BlogStorage{
		data: make(map[uuid.UUID]*admin.Blog),
	}
	return &s
}

func (s *BlogStorage) Create(blog *admin.Blog) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[blog.ID()]
	if ok {
		return storage.ErrConflict
	}

	s.data[blog.ID()] = blog
	return nil
}

func (s *BlogStorage) Read(id uuid.UUID) (*admin.Blog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	blog, ok := s.data[id]
	if !ok {
		return nil, storage.ErrNotFound
	}

	return blog, nil
}

func (s *BlogStorage) ReadByFeedURL(feedURL string) (*admin.Blog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, blog := range s.data {
		if blog.FeedURL() == feedURL {
			return blog, nil
		}
	}

	return nil, storage.ErrNotFound
}

func (s *BlogStorage) List(limit, offset int) ([]*admin.Blog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var blogs []*admin.Blog
	for _, blog := range s.data {
		blogs = append(blogs, blog)
	}

	start := offset
	end := min(offset+limit, len(blogs))
	if start >= len(blogs) {
		return nil, nil
	}

	return blogs[start:end], nil
}

func (s *BlogStorage) Update(blog *admin.Blog) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[blog.ID()]
	if !ok {
		return storage.ErrNotFound
	}

	s.data[blog.ID()] = blog
	return nil
}

func (s *BlogStorage) Delete(blog *admin.Blog) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[blog.ID()]
	if !ok {
		return storage.ErrNotFound
	}

	delete(s.data, blog.ID())
	return nil
}
