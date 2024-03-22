package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

// ensure PostStorage interface is satisfied
var _ PostStorage = (*PostgresPostStorage)(nil)

type PostStorage interface {
	Create(post admin.Post) error
	Read(id uuid.UUID) (admin.Post, error)
	List(limit, offset int) ([]admin.Post, error)
	ListByBlog(blog admin.Blog, limit, offset int) ([]admin.Post, error)
}

type dbPost struct {
	ID          uuid.UUID `db:"id"`
	BlogID      uuid.UUID `db:"blog_id"`
	URL         string    `db:"url"`
	Title       string    `db:"title"`
	Content     string    `db:"content"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type PostgresPostStorage struct {
	conn postgres.Conn
}

func NewPostgresPostStorage(conn postgres.Conn) *PostgresPostStorage {
	s := PostgresPostStorage{
		conn: conn,
	}
	return &s
}

func (s *PostgresPostStorage) marshal(post admin.Post) (dbPost, error) {
	row := dbPost{
		ID:          post.ID,
		BlogID:      post.BlogID,
		URL:         post.URL,
		Title:       post.Title,
		Content:     post.Content,
		PublishedAt: post.PublishedAt,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
	return row, nil
}

func (s *PostgresPostStorage) unmarshal(row dbPost) (admin.Post, error) {
	post := admin.Post{
		ID:          row.ID,
		BlogID:      row.BlogID,
		URL:         row.URL,
		Title:       row.Title,
		Content:     row.Content,
		PublishedAt: row.PublishedAt,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
	return post, nil
}

func (s *PostgresPostStorage) Create(post admin.Post) error {
	stmt := `
		INSERT INTO post
			(id, blog_id, url, title, content, published_at, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)`

	row, err := s.marshal(post)
	if err != nil {
		return err
	}

	args := []interface{}{
		row.ID,
		row.BlogID,
		row.URL,
		row.Title,
		row.Content,
		row.PublishedAt,
		row.CreatedAt,
		row.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	_, err = s.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return checkCreateError(err)
	}

	return nil
}

func (s *PostgresPostStorage) Read(id uuid.UUID) (admin.Post, error) {
	stmt := `
		SELECT
			id,
			blog_id,
			url,
			title,
			content,
			published_at,
			created_at,
			updated_at
		FROM post
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, id)
	if err != nil {
		return admin.Post{}, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbPost])
	if err != nil {
		return admin.Post{}, checkReadError(err)
	}

	return s.unmarshal(row)
}

func (s *PostgresPostStorage) List(limit, offset int) ([]admin.Post, error) {
	stmt := `
		SELECT
			id,
			blog_id,
			url,
			title,
			content,
			published_at,
			created_at,
			updated_at
		FROM post
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}

	postRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbPost])
	if err != nil {
		return nil, checkListError(err)
	}

	var posts []admin.Post
	for _, row := range postRows {
		post, err := s.unmarshal(row)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostgresPostStorage) ListByBlog(blog admin.Blog, limit, offset int) ([]admin.Post, error) {
	stmt := `
		SELECT
			id,
			blog_id,
			url,
			title,
			content,
			published_at,
			created_at,
			updated_at
		FROM post
		WHERE post.blog_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, blog.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	postRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbPost])
	if err != nil {
		return nil, checkListError(err)
	}

	var posts []admin.Post
	for _, row := range postRows {
		post, err := s.unmarshal(row)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
