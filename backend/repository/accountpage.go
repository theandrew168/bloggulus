package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

type AccountPageRepository struct {
	conn postgres.Conn
}

func NewAccountPageRepository(conn postgres.Conn) *AccountPageRepository {
	r := AccountPageRepository{
		conn: conn,
	}
	return &r
}

func (r *AccountPageRepository) Create(account *model.Account, page *model.Page) error {
	stmt := `
		INSERT INTO account_page
			(account_id, page_id, created_at, updated_at)
		VALUES
			($1, $2, $3, $4)`

	now := timeutil.Now()
	args := []any{
		account.ID(),
		page.ID(),
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

func (r *AccountPageRepository) Delete(account *model.Account, page *model.Page) error {
	stmt := `
		DELETE FROM account_page
		WHERE account_id = $1
			AND page_id = $2
		RETURNING account_id`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	args := []any{
		account.ID(),
		page.ID(),
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
