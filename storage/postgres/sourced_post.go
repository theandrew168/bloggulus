package postgres

import (
	"context"

	"github.com/theandrew168/bloggulus/models"
	"github.com/theandrew168/bloggulus/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type sourcedPostStorage struct {
	db *pgxpool.Pool
}

func NewSourcedPostStorage(db *pgxpool.Pool) storage.SourcedPost {
	return &sourcedPostStorage{
		db: db,
	}
}

func (s *sourcedPostStorage) ReadRecent(ctx context.Context, n int) ([]*models.SourcedPost, error) {
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

	var posts []*models.SourcedPost
	for rows.Next() {
		var post models.SourcedPost
		err := rows.Scan(&post.URL, &post.Title, &post.Updated, &post.BlogTitle)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (s *sourcedPostStorage) ReadRecentForUser(ctx context.Context, accountID int, n int) ([]*models.SourcedPost, error) {
	query := `
		SELECT
			post.url,
			post.title,
			post.updated,
			blog.title
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		INNER JOIN account_blog
			ON account_blog.blog_id = blog.blog_id
		WHERE account_blog.account_id = $1
		ORDER BY post.updated DESC
		LIMIT $2`
	rows, err := s.db.Query(ctx, query, accountID, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.SourcedPost
	for rows.Next() {
		var post models.SourcedPost
		err := rows.Scan(&post.URL, &post.Title, &post.Updated, &post.BlogTitle)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}
