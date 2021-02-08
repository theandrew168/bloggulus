package query

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Post struct {
	URL       string
	Title     string
	Updated   time.Time
	BlogTitle string
}

type PostStorage struct {
	db *pgxpool.Pool
}

func (s *PostStorage) ReadRecent(ctx context.Context, n int) ([]*Post, error) {
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

    var posts []*Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.URL, &post.Title, &post.Updated, &post.BlogTitle)
        if err != nil {
            return nil, err
        }

        posts = append(posts, &post)
    }

    return posts, nil
}

func (s *PostStorage) ReadRecentForUser(ctx context.Context, accountID int, n int) ([]*Post, error) {
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

    var posts []*Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.URL, &post.Title, &post.Updated, &post.BlogTitle)
        if err != nil {
            return nil, err
        }

        posts = append(posts, &post)
    }

    return posts, nil
}
