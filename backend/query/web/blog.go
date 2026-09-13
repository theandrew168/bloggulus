package webquery

import (
	"context"
	"strings"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/value"
)

// TODO: Should these use proper values like Name and URL?

type Blog struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	SiteURL     string    `db:"site_url"`
	IsFollowing bool      `db:"is_following"`
}

type BlogDetails struct {
	ID       uuid.UUID `db:"id"`
	FeedURL  string    `db:"feed_url"`
	SiteURL  string    `db:"site_url"`
	Title    string    `db:"title"`
	SyncedAt time.Time `db:"synced_at"`
	IsPublic bool      `db:"is_public"`
}

type BlogQuery struct {
	conn postgres.Conn
}

func NewBlog(conn postgres.Conn) *BlogQuery {
	qry := BlogQuery{
		conn: conn,
	}
	return &qry
}

// Powers the add / follow blogs page (admins only).
func (qry *BlogQuery) ListAll(ctx context.Context, accountID uuid.UUID) ([]Blog, error) {
	stmt := `
		-- name: WebQuery_Blog_ListAll
		SELECT
			blog.id,
			blog.title,
			blog.site_url,
			account_blog IS NOT NULL AS is_following
		FROM blog
		LEFT JOIN account_blog
			ON account_blog.blog_id = blog.id
			AND account_blog.account_id = $1
		ORDER BY blog.title ASC;
	`

	rows, err := qry.conn.Query(ctx, strings.TrimSpace(stmt), accountID)
	if err != nil {
		return nil, err
	}

	blogs, err := pgx.CollectRows(rows, pgx.RowToStructByName[Blog])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	return blogs, nil
}

// Powers the add / follow blogs page (non-admins, public and / or followed only).
func (qry *BlogQuery) ListVisible(ctx context.Context, accountID uuid.UUID) ([]Blog, error) {
	stmt := `
		-- name: WebQuery_Blog_ListVisible
		SELECT
			blog.id,
			blog.title,
			blog.site_url,
			account_blog IS NOT NULL AS is_following
		FROM blog
		LEFT JOIN account_blog
			ON account_blog.blog_id = blog.id
			AND account_blog.account_id = $1
		WHERE (blog.is_public = TRUE OR account_blog.blog_id IS NOT NULL)
		ORDER BY blog.title ASC;
	`

	rows, err := qry.conn.Query(ctx, strings.TrimSpace(stmt), accountID)
	if err != nil {
		return nil, err
	}

	blogs, err := pgx.CollectRows(rows, pgx.RowToStructByName[Blog])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	return blogs, nil
}

// Powers the blog details page (admin only).
func (qry *BlogQuery) ReadDetailsByID(ctx context.Context, blogID uuid.UUID) (BlogDetails, error) {
	stmt := `
		-- name: WebQuery_Blog_ReadDetailsByID
		SELECT
			blog.id,
			blog.feed_url,
			blog.site_url,
			blog.title,
			blog.synced_at,
			blog.is_public
		FROM blog
		WHERE blog.id = $1;
	`

	rows, err := qry.conn.Query(ctx, strings.TrimSpace(stmt), blogID)
	if err != nil {
		return BlogDetails{}, err
	}

	details, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[BlogDetails])
	if err != nil {
		return BlogDetails{}, postgres.CheckReadError(err)
	}

	return details, nil
}

// Powers the add / follow blogs page.
func (qry *BlogQuery) ReadDetailsByFeedURL(ctx context.Context, feedURL value.URL) (BlogDetails, error) {
	stmt := `
		-- name: WebQuery_Blog_ReadDetailsByFeedURL
		SELECT
			blog.id,
			blog.feed_url,
			blog.site_url,
			blog.title,
			blog.synced_at,
			blog.is_public
		FROM blog
		WHERE blog.feed_url = $1;
	`

	rows, err := qry.conn.Query(ctx, strings.TrimSpace(stmt), feedURL.Value())
	if err != nil {
		return BlogDetails{}, err
	}

	details, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[BlogDetails])
	if err != nil {
		return BlogDetails{}, postgres.CheckReadError(err)
	}

	return details, nil
}
