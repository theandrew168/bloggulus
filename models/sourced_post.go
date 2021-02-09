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
			post.url,
			post.title,
			post.updated,
			blog.title
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		ORDER BY post.updated DESC
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
			post.url,
			post.title,
			post.updated,
			blog.title
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		INNER JOIN follows
			ON follows.blog_id = blog.blog_id
		WHERE follows.account_id = $1
		ORDER BY post.updated DESC
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
