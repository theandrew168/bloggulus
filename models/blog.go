package models

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Blog struct {
	BlogID  int
	FeedURL string
	SiteURL string
	Title   string
}

type BlogStorage struct {
	db *pgxpool.Pool
}

func NewBlogStorage(db *pgxpool.Pool) *BlogStorage {
	return &BlogStorage{
		db: db,
	}
}

func (s *BlogStorage) Create(ctx context.Context, feedURL, siteURL, title string) (*Blog, error) {
	stmt := "INSERT INTO blog (feed_url, site_url, title) VALUES ($1, $2, $3) RETURNING blog_id"
	row := s.db.QueryRow(ctx, stmt, feedURL, siteURL, title)

	var blogID int
	err := row.Scan(&blogID)
	if err != nil {
		return nil, err
	}

	blog := &Blog{
		BlogID:  blogID,
		FeedURL: feedURL,
		SiteURL: siteURL,
		Title:   title,
	}

	return blog, nil
}

func (s *BlogStorage) Read(ctx context.Context, blogID int) (*Blog, error) {
	query := "SELECT * FROM blog WHERE blog_id = $1"
	row := s.db.QueryRow(ctx, query, blogID)

	var blog Blog
	err := row.Scan(&blog.BlogID, &blog.FeedURL, &blog.SiteURL, &blog.Title)
	if err != nil {
		return nil, err
	}

	return &blog, nil
}

func (s *BlogStorage) ReadAll(ctx context.Context) ([]*Blog, error) {
	query := "SELECT * FROM blog"
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []*Blog
	for rows.Next() {
		var blog Blog
		err := rows.Scan(&blog.BlogID, &blog.FeedURL, &blog.SiteURL, &blog.Title)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, &blog)
	}

	return blogs, nil
}
