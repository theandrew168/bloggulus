package repository

import (
	"context"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

type dbBlog struct {
	ID            uuid.UUID `db:"id"`
	FeedURL       string    `db:"feed_url"`
	SiteURL       string    `db:"site_url"`
	Title         string    `db:"title"`
	IsPublic      bool      `db:"is_public"`
	ETag          string    `db:"etag"`
	LastModified  string    `db:"last_modified"`
	SyncedAt      time.Time `db:"synced_at"`
	MetaCreatedAt time.Time `db:"meta_created_at"`
	MetaUpdatedAt time.Time `db:"meta_updated_at"`
}

func marshalBlog(blog *model.Blog) (dbBlog, error) {
	b := dbBlog{
		ID:            blog.ID(),
		FeedURL:       blog.FeedURL(),
		SiteURL:       blog.SiteURL(),
		Title:         blog.Title(),
		IsPublic:      blog.IsPublic(),
		ETag:          blog.ETag(),
		LastModified:  blog.LastModified(),
		SyncedAt:      blog.SyncedAt(),
		MetaCreatedAt: blog.Meta().CreatedAt(),
		MetaUpdatedAt: blog.Meta().UpdatedAt(),
	}
	return b, nil
}

func (b dbBlog) unmarshal() (*model.Blog, error) {
	metaParams := model.LoadMetaParams{
		CreatedAt: b.MetaCreatedAt,
		UpdatedAt: b.MetaUpdatedAt,
	}

	params := model.LoadBlogParams{
		ID:           b.ID,
		FeedURL:      b.FeedURL,
		SiteURL:      b.SiteURL,
		Title:        b.Title,
		IsPublic:     b.IsPublic,
		SyncedAt:     b.SyncedAt,
		ETag:         b.ETag,
		LastModified: b.LastModified,
		Meta:         model.LoadMeta(metaParams),
	}
	blog := model.LoadBlog(params)
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
			(id, feed_url, site_url, title, is_public, etag, last_modified, synced_at, meta_created_at, meta_updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	row, err := marshalBlog(blog)
	if err != nil {
		return err
	}

	args := []any{
		row.ID,
		row.FeedURL,
		row.SiteURL,
		row.Title,
		row.IsPublic,
		row.ETag,
		row.LastModified,
		row.SyncedAt,
		row.MetaCreatedAt,
		row.MetaUpdatedAt,
	}

	_, err = r.conn.Exec(context.Background(), stmt, args...)
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
			blog.is_public,
			blog.etag,
			blog.last_modified,
			blog.synced_at,
			blog.meta_created_at,
			blog.meta_updated_at
		FROM blog
		WHERE id = $1`

	rows, err := r.conn.Query(context.Background(), stmt, id)
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
			blog.is_public,
			blog.etag,
			blog.last_modified,
			blog.synced_at,
			blog.meta_created_at,
			blog.meta_updated_at
		FROM blog
		WHERE blog.feed_url = $1`

	rows, err := r.conn.Query(context.Background(), stmt, feedURL)
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
			blog.is_public,
			blog.etag,
			blog.last_modified,
			blog.synced_at,
			blog.meta_created_at,
			blog.meta_updated_at
		FROM blog
		ORDER BY blog.meta_created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.conn.Query(context.Background(), stmt, limit, offset)
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
			blog.is_public,
			blog.etag,
			blog.last_modified,
			blog.synced_at,
			blog.meta_created_at,
			blog.meta_updated_at
		FROM blog
		ORDER BY blog.meta_created_at DESC`

	rows, err := r.conn.Query(context.Background(), stmt)
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

	rows, err := r.conn.Query(context.Background(), stmt)
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
			feed_url = $2,
			site_url = $3,
			title = $4,
			is_public = $5,
			etag = $6,
			last_modified = $7,
			synced_at = $8,
			meta_updated_at = $9
		WHERE id = $1
			AND meta_updated_at = $10
		RETURNING meta_updated_at`

	row, err := marshalBlog(blog)
	if err != nil {
		return err
	}

	args := []any{
		row.ID,
		row.FeedURL,
		row.SiteURL,
		row.Title,
		row.IsPublic,
		row.ETag,
		row.LastModified,
		row.SyncedAt,
		now,
		row.MetaUpdatedAt,
	}

	rows, err := r.conn.Query(context.Background(), stmt, args...)
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[time.Time])
	if err != nil {
		return postgres.CheckUpdateError(err)
	}

	blog.Meta().Update(now)
	return nil
}

func (r *BlogRepository) Delete(blog *model.Blog) error {
	stmt := `
		DELETE FROM blog
		WHERE id = $1
		RETURNING id`

	rows, err := r.conn.Query(context.Background(), stmt, blog.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}
