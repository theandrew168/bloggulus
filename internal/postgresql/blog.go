package postgresql

import (
	"context"
	"errors"

	"github.com/theandrew168/bloggulus/internal/core"
)

func (s *storage) CreateBlog(ctx context.Context, blog *core.Blog) error {
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
	row := s.conn.QueryRow(ctx, stmt, args...)

	err := scan(row, &blog.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.CreateBlog(ctx, blog)
		}
		return err
	}

	return nil
}

func (s *storage) ReadBlog(ctx context.Context, id int) (core.Blog, error) {
	stmt := `
		SELECT
			id,
			feed_url,
			site_url,
			title
		FROM blog
		WHERE id = $1`
	row := s.conn.QueryRow(ctx, stmt, id)

	var blog core.Blog
	err := scan(
		row,
		&blog.ID,
		&blog.FeedURL,
		&blog.SiteURL,
		&blog.Title,
	)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.ReadBlog(ctx, id)
		}
		return core.Blog{}, err
	}

	return blog, nil
}

func (s *storage) ReadBlogs(ctx context.Context, limit, offset int) ([]core.Blog, error) {
	stmt := `
		SELECT
			id,
			feed_url,
			site_url,
			title
		FROM blog
		ORDER BY title ASC
		LIMIT $1 OFFSET $2`
	rows, err := s.conn.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// use make here to JSON encode as an empty array instead of null
	blogs := make([]core.Blog, 0)
	for rows.Next() {
		var blog core.Blog
		err := scan(
			rows,
			&blog.ID,
			&blog.FeedURL,
			&blog.SiteURL,
			&blog.Title,
		)
		if err != nil {
			if errors.Is(err, core.ErrRetry) {
				return s.ReadBlogs(ctx, limit, offset)
			}
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}
