package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

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

	err := row.Scan(&blog.ID)
	if err != nil {
		// https://github.com/jackc/pgx/wiki/Error-Handling
		// https://github.com/jackc/pgx/issues/474
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return core.ErrExist
			}
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
	err := row.Scan(
		&blog.ID,
		&blog.FeedURL,
		&blog.SiteURL,
		&blog.Title,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.Blog{}, core.ErrNotExist
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
		err := rows.Scan(
			&blog.ID,
			&blog.FeedURL,
			&blog.SiteURL,
			&blog.Title,
		)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}
