package storage

import (
	"context"

	"github.com/theandrew168/bloggulus/backend/database"
	"github.com/theandrew168/bloggulus/backend/domain"
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

func (s *Blog) Create(blog *domain.Blog) error {
	stmt := `
		INSERT INTO blog
			(feed_url, site_url, title, etag, last_modified)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id`

	args := []interface{}{
		blog.FeedURL,
		blog.SiteURL,
		blog.Title,
		blog.ETag,
		blog.LastModified,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &blog.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Blog) Read(id int) (domain.Blog, error) {
	stmt := `
		SELECT
			id,
			feed_url,
			site_url,
			title,
			etag,
			last_modified
		FROM blog
		WHERE id = $1`

	var blog domain.Blog
	dest := []interface{}{
		&blog.ID,
		&blog.FeedURL,
		&blog.SiteURL,
		&blog.Title,
		&blog.ETag,
		&blog.LastModified,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		return domain.Blog{}, err
	}

	return blog, nil
}

func (s *Blog) ReadAll(limit, offset int) ([]domain.Blog, error) {
	stmt := `
		SELECT
			id,
			feed_url,
			site_url,
			title,
			etag,
			last_modified
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
	blogs := make([]domain.Blog, 0)
	for rows.Next() {
		var blog domain.Blog
		dest := []interface{}{
			&blog.ID,
			&blog.FeedURL,
			&blog.SiteURL,
			&blog.Title,
			&blog.ETag,
			&blog.LastModified,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return blogs, nil
}

func (s *Blog) Update(blog domain.Blog) error {
	stmt := `
		UPDATE blog
		SET
			feed_url = $2,
			site_url = $3,
			title = $4,
			etag = $5,
			last_modified = $6
		WHERE id = $1`

	args := []interface{}{
		blog.ID,
		blog.FeedURL,
		blog.SiteURL,
		blog.Title,
		blog.ETag,
		blog.LastModified,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := s.db.Exec(ctx, stmt, args...)
	if err != nil {
		return err
	}

	return nil
}
