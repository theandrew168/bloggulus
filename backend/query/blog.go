package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

type BlogForAccount struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	SiteURL     string    `db:"site_url"`
	IsFollowing bool      `db:"is_following"`
}

// TODO: Paginate this (will need to add a CountBlogsForAccount method).
func (qry *Query) ListBlogsForAccount(account *model.Account) ([]BlogForAccount, error) {
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
		ORDER BY blog.created_at DESC`

	rows, err := qry.conn.Query(context.Background(), stmt, account.ID())
	if err != nil {
		return nil, err
	}

	blogs, err := pgx.CollectRows(rows, pgx.RowToStructByName[BlogForAccount])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	return blogs, nil
}
