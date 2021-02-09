package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type SourcedPost struct {
	URL       string
	Title     string
	Updated   time.Time
	BlogTitle string
}

type SourcedPostStorage struct {
	db *pgxpool.Pool
}

func NewSourcedPostStorage(db *pgxpool.Pool) *SourcedPostStorage {
	return &SourcedPostStorage{
		db: db,
	}
}

func (s *SourcedPostStorage) ReadRecent(ctx context.Context, n int) ([]*SourcedPost, error) {
	query := `
		SELECT
			posts.url,
			posts.title,
			posts.updated,
			blogs.title
		FROM posts
		INNER JOIN blogs
			ON blogs.blog_id = posts.blog_id
		ORDER BY posts.updated DESC
		LIMIT $1`
	rows, err := s.db.Query(ctx, query, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*SourcedPost
	for rows.Next() {
		var post SourcedPost
		err := rows.Scan(&post.URL, &post.Title, &post.Updated, &post.BlogTitle)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (s *SourcedPostStorage) ReadRecentForUser(ctx context.Context, accountID int, n int) ([]*SourcedPost, error) {
	query := `
		SELECT
			posts.url,
			posts.title,
			posts.updated,
			blogs.title
		FROM posts
		INNER JOIN blogs
			ON blogs.blog_id = posts.blog_id
		INNER JOIN follows
			ON follows.blog_id = blogs.blog_id
		WHERE follows.account_id = $1
		ORDER BY posts.updated DESC
		LIMIT $2`
	rows, err := s.db.Query(ctx, query, accountID, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*SourcedPost
	for rows.Next() {
		var post SourcedPost
		err := rows.Scan(&post.URL, &post.Title, &post.Updated, &post.BlogTitle)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}
