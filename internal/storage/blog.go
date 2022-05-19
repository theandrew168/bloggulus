package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/bloggulus"
	"github.com/theandrew168/bloggulus/internal/database"
)

type Blog struct {
	db database.Conn
}

func NewBlog(db database.Conn) *Blog {
	s := Blog{
		db: db,
	}
	return &s
}

func (s *Blog) Create(blog *bloggulus.Blog) error {
	stmt := `
		INSERT INTO blog
			(feed_url, site_url, title)
		VALUES
			($1, $2, $3)
		RETURNING id`

	args := []interface{}{
		blog.FeedURL,
		blog.SiteURL,
		blog.Title,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &blog.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Create(blog)
		}
		return err
	}

	return nil
}

func (s *Blog) Read(id int) (bloggulus.Blog, error) {
	stmt := `
		SELECT
			id,
			feed_url,
			site_url,
			title
		FROM blog
		WHERE id = $1`

	var blog bloggulus.Blog
	dest := []interface{}{
		&blog.ID,
		&blog.FeedURL,
		&blog.SiteURL,
		&blog.Title,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Read(id)
		}
		return bloggulus.Blog{}, err
	}

	return blog, nil
}

func (s *Blog) ReadAll(limit, offset int) ([]bloggulus.Blog, error) {
	stmt := `
		SELECT
			id,
			feed_url,
			site_url,
			title
		FROM blog
		ORDER BY title ASC
		LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := s.db.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// use "make" here to encode JSON as an empty array instead of null
	blogs := make([]bloggulus.Blog, 0)
	for rows.Next() {
		var blog bloggulus.Blog
		dest := []interface{}{
			&blog.ID,
			&blog.FeedURL,
			&blog.SiteURL,
			&blog.Title,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
				return s.ReadAll(limit, offset)
			}
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return blogs, nil
}
