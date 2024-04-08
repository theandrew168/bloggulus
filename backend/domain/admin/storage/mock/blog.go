package mock

import (
	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
)

// ensure BlogStorage interface is satisfied
var _ storage.BlogStorage = (*BlogStorage)(nil)

type BlogStorage struct {
	data map[uuid.UUID]*admin.Blog
}

func NewBlogStorage() *BlogStorage {
	s := BlogStorage{
		data: make(map[uuid.UUID]*admin.Blog),
	}
	return &s
}

func (s *BlogStorage) Create(blog *admin.Blog) error {
	_, ok := s.data[blog.ID()]
	if ok {
		return storage.ErrConflict
	}

	s.data[blog.ID()] = blog
	return nil
}

func (s *BlogStorage) Read(id uuid.UUID) (*admin.Blog, error) {
	blog, ok := s.data[id]
	if !ok {
		return nil, storage.ErrNotFound
	}

	return blog, nil
}

func (s *BlogStorage) ReadByFeedURL(feedURL string) (*admin.Blog, error) {
	for _, blog := range s.data {
		if blog.FeedURL() == feedURL {
			return blog, nil
		}
	}

	return nil, storage.ErrNotFound
}

func (s *BlogStorage) List(limit, offset int) ([]*admin.Blog, error) {
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

func (s *BlogStorage) Update(blog *admin.Blog) error {
	_, ok := s.data[blog.ID()]
	if !ok {
		return storage.ErrNotFound
	}

	s.data[blog.ID()] = blog
	return nil
}

func (s *BlogStorage) Delete(blog *admin.Blog) error {
	_, ok := s.data[blog.ID()]
	if !ok {
		return storage.ErrNotFound
	}

	delete(s.data, blog.ID())
	return nil
}
