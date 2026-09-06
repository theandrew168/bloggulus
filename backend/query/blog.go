package query

import (
	"context"
	"uuid"

	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/postgres"
)

type BlogForAccount struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	SiteURL     string    `db:"site_url"`
	IsFollowing bool      `db:"is_following"`
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

func (qry *BlogQuery) ListBlogsForAccount(accountID uuid.UUID) ([]BlogForAccount, error) {
	stmt := `
		SELECT
			blog.id,
			blog.title,
			blog.site_url,
			account_blog IS NOT NULL AS is_following
		FROM blog
		LEFT JOIN account_blog
			ON account_blog.blog_id = blog.id
			AND account_blog.account_id = $1
		ORDER BY blog.meta_created_at DESC;
	`

	rows, err := qry.conn.Query(context.Background(), stmt, accountID)
	if err != nil {
		return nil, err
	}

	blogs, err := pgx.CollectRows(rows, pgx.RowToStructByName[BlogForAccount])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	return blogs, nil
}
