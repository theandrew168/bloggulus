package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

// ensure BlogStorage interface is satisfied
var _ storage.BlogStorage = (*PostgresBlogStorage)(nil)

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
	conn postgres.Conn
}

func NewPostgresBlogStorage(conn postgres.Conn) *PostgresBlogStorage {
	s := PostgresBlogStorage{
		conn: conn,
	}
	return &s
}

func (s *PostgresBlogStorage) marshal(blog admin.Blog) (dbBlog, error) {
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

func (s *PostgresBlogStorage) unmarshal(row dbBlog) (admin.Blog, error) {
	blog := admin.Blog{
		ID:           row.ID,
		FeedURL:      row.FeedURL,
		SiteURL:      row.SiteURL,
		Title:        row.Title,
		ETag:         row.ETag,
		LastModified: row.LastModified,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	return blog, nil
}

func (s *PostgresBlogStorage) Create(blog admin.Blog) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	_, err = s.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return postgres.CheckCreateError(err)
	}

	return nil
}

func (s *PostgresBlogStorage) Read(id uuid.UUID) (admin.Blog, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, id)
	if err != nil {
		return admin.Blog{}, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbBlog])
	if err != nil {
		return admin.Blog{}, postgres.CheckReadError(err)
	}

	return s.unmarshal(row)
}

func (s *PostgresBlogStorage) List(limit, offset int) ([]admin.Blog, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}

	blogRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbBlog])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	var blogs []admin.Blog
	for _, row := range blogRows {
		blog, err := s.unmarshal(row)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func (s *PostgresBlogStorage) Update(blog admin.Blog) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, args...)
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[time.Time])
	if err != nil {
		return postgres.CheckUpdateError(err)
	}

	blog.UpdatedAt = now
	return nil
}
