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

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	_, err := r.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return postgres.CheckCreateError(err)
	}

	return nil
}

func (r *AccountBlogRepository) Count(account *model.Account, blog *model.Blog) (int, error) {
	stmt := `
		SELECT count(*)
		FROM account_blog
		WHERE account_id = $1
			AND blog_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	args := []any{
		account.ID(),
		blog.ID(),
	}

	rows, err := r.conn.Query(ctx, stmt, args...)
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}

func (r *AccountBlogRepository) Delete(account *model.Account, blog *model.Blog) error {
	stmt := `
		DELETE FROM account_blog
		WHERE account_id = $1
			AND blog_id = $2
		RETURNING account_id`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	args := []any{
		account.ID(),
		blog.ID(),
	}

	rows, err := r.conn.Query(ctx, stmt, args...)
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}
