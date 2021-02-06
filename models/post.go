package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Post struct {
	PostID  int
	BlogID  int
	URL     string
	Title   string
	Updated time.Time
}

type PostStorage struct {
	db *pgxpool.Pool
}

func NewPostStorage(db *pgxpool.Pool) *PostStorage {
	return &PostStorage{
		db: db,
	}
}

func (s *PostStorage) Create(ctx context.Context, blogID int, URL, title string, updated time.Time) (*Post, error) {
	stmt := "INSERT INTO posts (blog_id, url, title, updated) VALUES ($1, $2, $3, $4) RETURNING post_id"
	row := s.db.QueryRow(ctx, stmt, blogID, URL, title, updated)

	var postID int
	err := row.Scan(&postID)
	if err != nil {
		return nil, err
	}

	post := &Post{
		PostID:  postID,
		BlogID:  blogID,
		URL:     URL,
		Title:   title,
		Updated: updated,
	}

	return post, nil
}

func (s *PostStorage) ReadRecent(ctx context.Context, n int) ([]*Post, error) {
	query := "SELECT * FROM posts ORDER BY updated DESC LIMIT $1"
	rows, err := s.db.Query(ctx, query, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.PostID, &post.BlogID, &post.URL, &post.Title, &post.Updated)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}
