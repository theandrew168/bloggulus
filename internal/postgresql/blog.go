package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/internal/core"
)

type blogStorage struct {
	conn *pgxpool.Pool
}

func NewBlogStorage(conn *pgxpool.Pool) core.BlogStorage {
	s := blogStorage{
		conn: conn,
	}
	return &s
}

func (s *blogStorage) Create(ctx context.Context, blog *core.Blog) error {
	stmt := `
		INSERT INTO blog
			(feed_url, site_url, title)
		VALUES
			($1, $2, $3)
		RETURNING blog_id`
	args := []interface{}{
		blog.FeedURL,
		blog.SiteURL,
		blog.Title,
	}
	row := s.conn.QueryRow(ctx, stmt, args...)

	err := row.Scan(&blog.BlogID)
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

func (s *blogStorage) ReadAll(ctx context.Context) ([]core.Blog, error) {
	stmt := "SELECT * FROM blog"
	rows, err := s.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []core.Blog
	for rows.Next() {
		var blog core.Blog
		err := rows.Scan(
			&blog.BlogID,
			&blog.FeedURL,
			&blog.SiteURL,
			&blog.Title)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}
