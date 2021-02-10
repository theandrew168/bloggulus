package postgres

import (
	"context"

	"github.com/theandrew168/bloggulus/models"
	"github.com/theandrew168/bloggulus/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type postStorage struct {
	db *pgxpool.Pool
}

func NewPostStorage(db *pgxpool.Pool) storage.Post {
	return &postStorage{
		db: db,
	}
}

func (s *postStorage) Create(ctx context.Context, post *models.Post) (*models.Post, error) {
	command := "INSERT INTO post (blog_id, url, title, updated) VALUES ($1, $2, $3, $4) RETURNING post_id"
	row := s.db.QueryRow(ctx, command, post.BlogID, post.URL, post.Title, post.Updated)

	err := row.Scan(&post.PostID)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postStorage) Read(ctx context.Context, postID int) (*models.Post, error) {
	query := "SELECT * FROM post WHERE post_id = $1"
	row := s.db.QueryRow(ctx, query, postID)

	var post models.Post
	err := row.Scan(&post.PostID, &post.BlogID, &post.URL, &post.Title, &post.Updated)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *postStorage) Delete(ctx context.Context, postID int) error {
	command := "DELETE FROM post WHERE post_id = $1"
	_, err := s.db.Exec(ctx, command, postID)
	return err
}
