package mock

import (
	"sort"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
)

// ensure PostStorage interface is satisfied
var _ storage.PostStorage = (*PostStorage)(nil)

type PostStorage struct {
	data map[uuid.UUID]*admin.Post
}

func NewPostStorage() *PostStorage {
	s := PostStorage{
		data: make(map[uuid.UUID]*admin.Post),
	}
	return &s
}

func (s *PostStorage) Create(post *admin.Post) error {
	_, ok := s.data[post.ID()]
	if ok {
		return storage.ErrConflict
	}

	s.data[post.ID()] = post
	return nil
}

func (s *PostStorage) Read(id uuid.UUID) (*admin.Post, error) {
	post, ok := s.data[id]
	if !ok {
		return nil, storage.ErrNotFound
	}

	return post, nil
}

func (s *PostStorage) ReadByURL(url string) (*admin.Post, error) {
	for _, post := range s.data {
		if post.URL() == url {
			return post, nil
		}
	}

	return nil, storage.ErrNotFound
}

func (s *PostStorage) List(limit, offset int) ([]*admin.Post, error) {
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

func (s *PostStorage) ListByBlog(blog *admin.Blog, limit, offset int) ([]*admin.Post, error) {
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

	sort.Slice(posts, func(i, j int) bool { return posts[i].PublishedAt().After(posts[j].PublishedAt()) })
	return posts[start:end], nil
}

func (s *PostStorage) Update(post *admin.Post) error {
	_, ok := s.data[post.ID()]
	if !ok {
		return storage.ErrNotFound
	}

	s.data[post.ID()] = post
	return nil
}

func (s *PostStorage) Delete(post *admin.Post) error {
	_, ok := s.data[post.ID()]
	if !ok {
		return storage.ErrNotFound
	}

	delete(s.data, post.ID())
	return nil
}
