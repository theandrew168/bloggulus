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
	db *pgxpool.Pool
}

func NewBlogStorage(db *pgxpool.Pool) core.BlogStorage {
	s := blogStorage{
		db: db,
	}
	return &s
}

func (s *blogStorage) Create(ctx context.Context, blog *core.Blog) (*core.Blog, error) {
	command := "INSERT INTO blog (feed_url, site_url, title) VALUES ($1, $2, $3) RETURNING blog_id"
	row := s.db.QueryRow(ctx, command, blog.FeedURL, blog.SiteURL, blog.Title)

	err := row.Scan(&blog.BlogID)
	if err != nil {
		// https://github.com/jackc/pgx/wiki/Error-Handling
		// https://github.com/jackc/pgx/issues/474
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, core.ErrExist
			}
		}
		return nil, err
	}

	return blog, nil
}

func (s *blogStorage) Read(ctx context.Context, blogID int) (*core.Blog, error) {
	query := "SELECT * FROM blog WHERE blog_id = $1"
	row := s.db.QueryRow(ctx, query, blogID)

	var blog core.Blog
	err := row.Scan(&blog.BlogID, &blog.FeedURL, &blog.SiteURL, &blog.Title)
	if err != nil {
		return nil, err
	}

	return &blog, nil
}

func (s *blogStorage) ReadByURL(ctx context.Context, feedURL string) (*core.Blog, error) {
	query := "SELECT * FROM blog WHERE feed_url = $1"
	row := s.db.QueryRow(ctx, query, feedURL)

	var blog core.Blog
	err := row.Scan(&blog.BlogID, &blog.FeedURL, &blog.SiteURL, &blog.Title)
	if err != nil {
		return nil, err
	}

	return &blog, nil
}

func (s *blogStorage) ReadAll(ctx context.Context) ([]*core.Blog, error) {
	query := "SELECT * FROM blog"
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []*core.Blog
	for rows.Next() {
		var blog core.Blog
		err := rows.Scan(&blog.BlogID, &blog.FeedURL, &blog.SiteURL, &blog.Title)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, &blog)
	}

	return blogs, nil
}

func (s *blogStorage) ReadFollowedForUser(ctx context.Context, accountID int) ([]*core.Blog, error) {
	query := `
		SELECT
			blog.*
		FROM blog
		INNER JOIN account_blog
			ON account_blog.blog_id = blog.blog_id
		WHERE account_blog.account_id = $1`
	rows, err := s.db.Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []*core.Blog
	for rows.Next() {
		var blog core.Blog
		err := rows.Scan(&blog.BlogID, &blog.FeedURL, &blog.SiteURL, &blog.Title)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, &blog)
	}

	return blogs, nil
}

func (s *blogStorage) ReadUnfollowedForUser(ctx context.Context, accountID int) ([]*core.Blog, error) {
	query := `
		SELECT
			blog.*
		FROM blog
		WHERE NOT EXISTS (
			SELECT 1
			FROM account_blog
			WHERE account_blog.blog_id = blog.blog_id
			AND account_blog.account_id = $1
		)`
	rows, err := s.db.Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []*core.Blog
	for rows.Next() {
		var blog core.Blog
		err := rows.Scan(&blog.BlogID, &blog.FeedURL, &blog.SiteURL, &blog.Title)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, &blog)
	}

	return blogs, nil
}

func (s *blogStorage) Delete(ctx context.Context, blogID int) error {
	command := "DELETE FROM blog WHERE blog_id = $1"
	_, err := s.db.Exec(ctx, command, blogID)
	return err
}
