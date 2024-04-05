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

// ensure BlogStorage interface is satisfied
var _ BlogStorage = (*MockBlogStorage)(nil)

type MockBlogStorage struct {
	data map[uuid.UUID]*admin.Blog
}

func NewMockBlogStorage() *MockBlogStorage {
	s := MockBlogStorage{
		data: make(map[uuid.UUID]*admin.Blog),
	}
	return &s
}

func (s *MockBlogStorage) Create(blog *admin.Blog) error {
	_, ok := s.data[blog.ID()]
	if ok {
		return ErrConflict
	}

	s.data[blog.ID()] = blog
	return nil
}

func (s *MockBlogStorage) Read(id uuid.UUID) (*admin.Blog, error) {
	blog, ok := s.data[id]
	if !ok {
		return nil, ErrNotFound
	}

	return blog, nil
}

func (s *MockBlogStorage) ReadByFeedURL(feedURL string) (*admin.Blog, error) {
	for _, blog := range s.data {
		if blog.FeedURL() == feedURL {
			return blog, nil
		}
	}

	return nil, ErrNotFound
}

func (s *MockBlogStorage) List(limit, offset int) ([]*admin.Blog, error) {
	var blogs []*admin.Blog
	for _, blog := range s.data {
		blogs = append(blogs, blog)
	}

	start := offset
	end := offset + limit
	if start >= len(blogs) || end >= len(blogs) {
		return nil, nil
	}

	return blogs[start:end], nil
}

func (s *MockBlogStorage) Update(blog *admin.Blog) error {
	_, ok := s.data[blog.ID()]
	if !ok {
		return ErrNotFound
	}

	s.data[blog.ID()] = blog
	return nil
}

func (s *MockBlogStorage) Delete(blog *admin.Blog) error {
	_, ok := s.data[blog.ID()]
	if !ok {
		return ErrNotFound
	}

	delete(s.data, blog.ID())
	return nil
}
