package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

type dbToken struct {
	ID        uuid.UUID `db:"id"`
	AccountID uuid.UUID `db:"account_id"`
	Hash      string    `db:"hash"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func marshalToken(token *model.Token) (dbToken, error) {
	t := dbToken{
		ID:        token.ID(),
		AccountID: token.AccountID(),
		Hash:      token.Hash(),
		ExpiresAt: token.ExpiresAt(),
		CreatedAt: token.CreatedAt(),
		UpdatedAt: token.UpdatedAt(),
	}
	return t, nil
}

func (t dbToken) unmarshal() (*model.Token, error) {
	token := model.LoadToken(
		t.ID,
		t.AccountID,
		t.Hash,
		t.ExpiresAt,
		t.CreatedAt,
		t.UpdatedAt,
	)
	return token, nil
}

type TokenStorage struct {
	conn postgres.Conn
}

func NewTokenStorage(conn postgres.Conn) *TokenStorage {
	s := TokenStorage{
		conn: conn,
	}
	return &s
}

func (s *TokenStorage) Create(token *model.Token) error {
	stmt := `
		INSERT INTO token
			(id, account_id, hash, expires_at, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6)`

	row, err := marshalToken(token)
	if err != nil {
		return err
	}

	args := []any{
		row.ID,
		row.AccountID,
		row.Hash,
		row.ExpiresAt,
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

func (s *TokenStorage) Read(id uuid.UUID) (*model.Token, error) {
	stmt := `
		SELECT
			token.id,
			token.account_id,
			token.hash,
			token.expires_at,
			token.created_at,
			token.updated_at
		FROM token
		WHERE token.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbToken])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (s *TokenStorage) Delete(token *model.Token) error {
	stmt := `
		DELETE FROM token
		WHERE id = $1
		RETURNING id`

	err := token.CheckDelete()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, token.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}
