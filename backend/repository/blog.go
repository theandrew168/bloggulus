package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

type dbBlog struct {
	ID           uuid.UUID `db:"id"`
	FeedURL      string    `db:"feed_url"`
	SiteURL      string    `db:"site_url"`
	Title        string    `db:"title"`
	ETag         string    `db:"etag"`
	LastModified string    `db:"last_modified"`
	SyncedAt     time.Time `db:"synced_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func marshalBlog(blog *model.Blog) (dbBlog, error) {
	b := dbBlog{
		ID:           blog.ID(),
		FeedURL:      blog.FeedURL(),
		SiteURL:      blog.SiteURL(),
		Title:        blog.Title(),
		ETag:         blog.ETag(),
		LastModified: blog.LastModified(),
		SyncedAt:     blog.SyncedAt(),
		CreatedAt:    blog.CreatedAt(),
		UpdatedAt:    blog.UpdatedAt(),
	}
	return b, nil
}

func (b dbBlog) unmarshal() (*model.Blog, error) {
	blog := model.LoadBlog(
		b.ID,
		b.FeedURL,
		b.SiteURL,
		b.Title,
		b.ETag,
		b.LastModified,
		b.SyncedAt,
		b.CreatedAt,
		b.UpdatedAt,
	)
	return blog, nil
}

type BlogRepository struct {
	conn postgres.Conn
}

func NewBlogRepository(conn postgres.Conn) *BlogRepository {
	r := BlogRepository{
		conn: conn,
	}
	return &r
}

func (r *BlogRepository) Create(blog *model.Blog) error {
	stmt := `
		INSERT INTO blog
			(id, feed_url, site_url, title, etag, last_modified, synced_at, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	row, err := marshalBlog(blog)
	if err != nil {
		return err
	}

	args := []any{
		row.ID,
		row.FeedURL,
		row.SiteURL,
		row.Title,
		row.ETag,
		row.LastModified,
		row.SyncedAt,
		row.CreatedAt,
		row.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	_, err = r.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return postgres.CheckCreateError(err)
	}

	return nil
}

func (r *BlogRepository) Read(id uuid.UUID) (*model.Blog, error) {
	stmt := `
		SELECT
			blog.id,
			blog.feed_url,
			blog.site_url,
			blog.title,
			blog.etag,
			blog.last_modified,
			blog.synced_at,
			blog.created_at,
			blog.updated_at
		FROM blog
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbBlog])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *BlogRepository) ReadByFeedURL(feedURL string) (*model.Blog, error) {
	stmt := `
		SELECT
			blog.id,
			blog.feed_url,
			blog.site_url,
			blog.title,
			blog.etag,
			blog.last_modified,
			blog.synced_at,
			blog.created_at,
			blog.updated_at
		FROM blog
		WHERE blog.feed_url = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, feedURL)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbBlog])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *BlogRepository) List(limit, offset int) ([]*model.Blog, error) {
	stmt := `
		SELECT
			blog.id,
			blog.feed_url,
			blog.site_url,
			blog.title,
			blog.etag,
			blog.last_modified,
			blog.synced_at,
			blog.created_at,
			blog.updated_at
		FROM blog
		ORDER BY blog.created_at DESC
		LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}

	blogRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbBlog])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	var blogs []*model.Blog
	for _, row := range blogRows {
		blog, err := row.unmarshal()
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}

// DEPRECATED
func (r *BlogRepository) ListAll() ([]*model.Blog, error) {
	stmt := `
		SELECT
			blog.id,
			blog.feed_url,
			blog.site_url,
			blog.title,
			blog.etag,
			blog.last_modified,
			blog.synced_at,
			blog.created_at,
			blog.updated_at
		FROM blog
		ORDER BY blog.created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	blogRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbBlog])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	var blogs []*model.Blog
	for _, row := range blogRows {
		blog, err := row.unmarshal()
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func (r *BlogRepository) Count() (int, error) {
	stmt := `
		SELECT count(*)
		FROM blog`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}

func (r *BlogRepository) Update(blog *model.Blog) error {
	now := timeutil.Now()
	stmt := `
		UPDATE blog
		SET
			feed_url = $1,
			site_url = $2,
			title = $3,
			etag = $4,
			last_modified = $5,
			synced_at = $6,
			updated_at = $7
		WHERE id = $8
			AND updated_at = $9
		RETURNING updated_at`

	row, err := marshalBlog(blog)
	if err != nil {
		return err
	}

	args := []any{
		row.FeedURL,
		row.SiteURL,
		row.Title,
		row.ETag,
		row.LastModified,
		row.SyncedAt,
		now,
		row.ID,
		row.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, args...)
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[time.Time])
	if err != nil {
		return postgres.CheckUpdateError(err)
	}

	blog.SetUpdatedAt(now)
	return nil
}

func (r *BlogRepository) Delete(blog *model.Blog) error {
	stmt := `
		DELETE FROM blog
		WHERE id = $1
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, blog.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}
