package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/database"
	"github.com/theandrew168/bloggulus/backend/domain"
)

// ensure BlogStorage interface is satisfied
var _ BlogStorage = (*PostgresBlogStorage)(nil)

type BlogStorage interface {
	Create(itinerary domain.Blog) error
	Read(id uuid.UUID) (domain.Blog, error)
	List(limit, offset int) ([]domain.Blog, error)
	Update(blog domain.Blog) error
}

type dbBlog struct {
	ID           uuid.UUID `db:"id"`
	FeedURL      string    `db:"feed_url"`
	SiteURL      string    `db:"site_url"`
	Title        string    `db:"title"`
	ETag         string    `db:"etag"`
	LastModified string    `db:"last_modified"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type PostgresBlogStorage struct {
	conn database.Conn
}

func NewPostgresBlogStorage(conn database.Conn) *PostgresBlogStorage {
	s := PostgresBlogStorage{
		conn: conn,
	}
	return &s
}

func (s *PostgresBlogStorage) marshal(blog domain.Blog) (dbBlog, error) {
	row := dbBlog{
		ID:           blog.ID,
		FeedURL:      blog.FeedURL,
		SiteURL:      blog.SiteURL,
		Title:        blog.Title,
		ETag:         blog.ETag,
		LastModified: blog.LastModified,
		CreatedAt:    blog.CreatedAt,
		UpdatedAt:    blog.UpdatedAt,
	}
	return row, nil
}

func (s *PostgresBlogStorage) unmarshal(row dbBlog) (domain.Blog, error) {
	blog := domain.LoadBlog(
		row.ID,
		row.FeedURL,
		row.SiteURL,
		row.Title,
		row.ETag,
		row.LastModified,
		row.CreatedAt,
		row.UpdatedAt,
	)
	return blog, nil
}

func (s *PostgresBlogStorage) Create(blog domain.Blog) error {
	stmt := `
		INSERT INTO blog
			(id, feed_url, site_url, title, etag, last_modified, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)`

	row, err := s.marshal(blog)
	if err != nil {
		return err
	}

	args := []interface{}{
		row.ID,
		row.FeedURL,
		row.SiteURL,
		row.Title,
		row.ETag,
		row.LastModified,
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

func (s *PostgresBlogStorage) Read(id uuid.UUID) (domain.Blog, error) {
	stmt := `
		SELECT
			id,
			feed_url,
			site_url,
			title,
			etag,
			last_modified,
			created_at,
			updated_at
		FROM blog
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, id)
	if err != nil {
		return domain.Blog{}, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbBlog])
	if err != nil {
		return domain.Blog{}, checkReadError(err)
	}

	return s.unmarshal(row)
}

func (s *PostgresBlogStorage) List(limit, offset int) ([]domain.Blog, error) {
	stmt := `
		SELECT
			id,
			feed_url,
			site_url,
			title,
			etag,
			last_modified,
			created_at,
			updated_at
		FROM blog
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}

	blogRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbBlog])
	if err != nil {
		return nil, checkListError(err)
	}

	var blogs []domain.Blog
	for _, row := range blogRows {
		blog, err := s.unmarshal(row)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func (s *PostgresBlogStorage) Update(blog domain.Blog) error {
	now := time.Now()
	stmt := `
		UPDATE blog
		SET
			feed_url = $1,
			site_url = $2,
			title = $3,
			etag = $4,
			last_modified = $5,
			updated_at = $6
		WHERE id = $7
		  AND updated_at = $6
		RETURNING updated_at`

	row, err := s.marshal(blog)
	if err != nil {
		return err
	}

	args := []any{
		row.FeedURL,
		row.SiteURL,
		row.Title,
		row.ETag,
		row.LastModified,
		now,
		row.ID,
		row.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, args...)
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[time.Time])
	if err != nil {
		return checkUpdateError(err)
	}

	blog.UpdatedAt = now
	return nil
}
