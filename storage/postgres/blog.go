package postgres

import (
	"context"

	"github.com/theandrew168/bloggulus/models"
	"github.com/theandrew168/bloggulus/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type blogStorage struct {
	db *pgxpool.Pool
}

func NewBlogStorage(db *pgxpool.Pool) storage.Blog {
	return &blogStorage{
		db: db,
	}
}

func (s *blogStorage) Create(ctx context.Context, blog *models.Blog) (*models.Blog, error) {
	command := "INSERT INTO blog (feed_url, site_url, title) VALUES ($1, $2, $3) RETURNING blog_id"
	row := s.db.QueryRow(ctx, command, blog.FeedURL, blog.SiteURL, blog.Title)

	err := row.Scan(&blog.BlogID)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *blogStorage) Read(ctx context.Context, blogID int) (*models.Blog, error) {
	query := "SELECT * FROM blog WHERE blog_id = $1"
	row := s.db.QueryRow(ctx, query, blogID)

	var blog models.Blog
	err := row.Scan(&blog.BlogID, &blog.FeedURL, &blog.SiteURL, &blog.Title)
	if err != nil {
		return nil, err
	}

	return &blog, nil
}

func (s *blogStorage) ReadAll(ctx context.Context) ([]*models.Blog, error) {
	query := "SELECT * FROM blog"
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []*models.Blog
	for rows.Next() {
		var blog models.Blog
		err := rows.Scan(&blog.BlogID, &blog.FeedURL, &blog.SiteURL, &blog.Title)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, &blog)
	}

	return blogs, nil
}

func (s *blogStorage) ReadFollowedForUser(ctx context.Context, accountID int) ([]*models.Blog, error) {
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

	var blogs []*models.Blog
	for rows.Next() {
		var blog models.Blog
		err := rows.Scan(&blog.BlogID, &blog.FeedURL, &blog.SiteURL, &blog.Title)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, &blog)
	}

	return blogs, nil
}

func (s *blogStorage) ReadUnfollowedForUser(ctx context.Context, accountID int) ([]*models.Blog, error) {
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

	var blogs []*models.Blog
	for rows.Next() {
		var blog models.Blog
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
