package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

type AccountBlogRepository struct {
	conn postgres.Conn
}

func NewAccountBlogRepository(conn postgres.Conn) *AccountBlogRepository {
	r := AccountBlogRepository{
		conn: conn,
	}
	return &r
}

func (r *AccountBlogRepository) Create(account *model.Account, blog *model.Blog) error {
	stmt := `
		INSERT INTO account_blog
			(account_id, blog_id, created_at, updated_at)
		VALUES
			($1, $2, $3, $4)`

	now := timeutil.Now()
	args := []any{
		account.ID(),
		blog.ID(),
		now,
		now,
	}

	_, err := r.conn.Exec(context.Background(), stmt, args...)
	if err != nil {
		return postgres.CheckCreateError(err)
	}

	return nil
}

func (r *AccountBlogRepository) Delete(account *model.Account, blog *model.Blog) error {
	stmt := `
		DELETE FROM account_blog
		WHERE account_id = $1
			AND blog_id = $2
		RETURNING account_id`

	args := []any{
		account.ID(),
		blog.ID(),
	}

	rows, err := r.conn.Query(context.Background(), stmt, args...)
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}
