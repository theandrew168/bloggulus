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
	command := "INSERT INTO post (blog_id, url, title, updated) VALUES ($1, $2, $3, $4) RETURNING post_id"
	row := s.db.QueryRow(ctx, command, blogID, URL, title, updated)

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

func (s *PostStorage) Read(ctx context.Context, postID int) (*Post, error) {
	query := "SELECT * FROM post WHERE post_id = $1"
	row := s.db.QueryRow(ctx, query, postID)

	var post Post
	err := row.Scan(&post.PostID, &post.BlogID, &post.URL, &post.Title, &post.Updated)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *PostStorage) Delete(ctx context.Context, postID int) error {
	command := "DELETE FROM post WHERE post_id = $1"
	_, err := s.db.Exec(ctx, command, postID)
	return err
}
