package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/database"
	"github.com/theandrew168/bloggulus/backend/domain"
)

// ensure PostStorage interface is satisfied
var _ PostStorage = (*PostgresPostStorage)(nil)

type PostStorage interface {
	Create(post domain.Post) error
	Read(id uuid.UUID) (domain.Post, error)
	List(limit, offset int) ([]domain.Post, error)
	ListByBlog(blog domain.Blog, limit, offset int) ([]domain.Post, error)
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
	conn database.Conn
}

func NewPostgresPostStorage(conn database.Conn) *PostgresPostStorage {
	s := PostgresPostStorage{
		conn: conn,
	}
	return &s
}

func (s *PostgresPostStorage) marshal(post domain.Post) (dbPost, error) {
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

func (s *PostgresPostStorage) unmarshal(row dbPost) (domain.Post, error) {
	post := domain.Post{
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

func (s *PostgresPostStorage) Create(post domain.Post) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	_, err = s.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return checkCreateError(err)
	}

	return nil
}

func (s *PostgresPostStorage) Read(id uuid.UUID) (domain.Post, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, id)
	if err != nil {
		return domain.Post{}, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbPost])
	if err != nil {
		return domain.Post{}, checkReadError(err)
	}

	return s.unmarshal(row)
}

func (s *PostgresPostStorage) List(limit, offset int) ([]domain.Post, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}

	postRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbPost])
	if err != nil {
		return nil, checkListError(err)
	}

	var posts []domain.Post
	for _, row := range postRows {
		post, err := s.unmarshal(row)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostgresPostStorage) ListByBlog(blog domain.Blog, limit, offset int) ([]domain.Post, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, blog.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	postRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbPost])
	if err != nil {
		return nil, checkListError(err)
	}

	var posts []domain.Post
	for _, row := range postRows {
		post, err := s.unmarshal(row)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
