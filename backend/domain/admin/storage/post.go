package storage

import (
	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
)

type PostStorage interface {
	Create(post *admin.Post) error
	Read(id uuid.UUID) (*admin.Post, error)
	ReadByURL(url string) (*admin.Post, error)
	List(limit, offset int) ([]*admin.Post, error)
	ListByBlog(blog *admin.Blog, limit, offset int) ([]*admin.Post, error)
	Update(post *admin.Post) error
	Delete(post *admin.Post) error
}

// ensure PostStorage interface is satisfied
var _ PostStorage = (*MockPostStorage)(nil)

type MockPostStorage struct {
	data map[uuid.UUID]*admin.Post
}

func NewMockPostStorage() *MockPostStorage {
	s := MockPostStorage{
		data: make(map[uuid.UUID]*admin.Post),
	}
	return &s
}

func (s *MockPostStorage) Create(post *admin.Post) error {
	_, ok := s.data[post.ID()]
	if ok {
		return ErrConflict
	}

	s.data[post.ID()] = post
	return nil
}

func (s *MockPostStorage) Read(id uuid.UUID) (*admin.Post, error) {
	post, ok := s.data[id]
	if !ok {
		return nil, ErrNotFound
	}

	return post, nil
}

func (s *MockPostStorage) ReadByURL(url string) (*admin.Post, error) {
	for _, post := range s.data {
		if post.URL() == url {
			return post, nil
		}
	}

	return nil, ErrNotFound
}

func (s *MockPostStorage) List(limit, offset int) ([]*admin.Post, error) {
	var posts []*admin.Post
	for _, post := range s.data {
		posts = append(posts, post)
	}

	start := offset
	end := offset + limit
	if start >= len(posts) || end >= len(posts) {
		return nil, nil
	}

	return posts[start:end], nil
}

func (s *MockPostStorage) ListByBlog(blog *admin.Blog, limit, offset int) ([]*admin.Post, error) {
	var posts []*admin.Post
	for _, post := range s.data {
		if post.BlogID() != blog.ID() {
			continue
		}

		posts = append(posts, post)
	}

	start := offset
	end := offset + limit
	if start >= len(posts) || end >= len(posts) {
		return nil, nil
	}

	return posts[start:end], nil
}

func (s *MockPostStorage) Update(post *admin.Post) error {
	_, ok := s.data[post.ID()]
	if !ok {
		return ErrNotFound
	}

	s.data[post.ID()] = post
	return nil
}

func (s *MockPostStorage) Delete(post *admin.Post) error {
	_, ok := s.data[post.ID()]
	if !ok {
		return ErrNotFound
	}

	delete(s.data, post.ID())
	return nil
}
