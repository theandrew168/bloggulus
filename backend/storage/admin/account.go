package admin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model/admin"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

type dbAccount struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	IsAdmin      bool      `db:"is_admin"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func marshalAccount(account *admin.Account) (dbAccount, error) {
	a := dbAccount{
		ID:           account.ID(),
		Username:     account.Username(),
		PasswordHash: account.PasswordHash(),
		IsAdmin:      account.IsAdmin(),
		CreatedAt:    account.CreatedAt(),
		UpdatedAt:    account.UpdatedAt(),
	}
	return a, nil
}

func (a dbAccount) unmarshal() (*admin.Account, error) {
	account := admin.LoadAccount(
		a.ID,
		a.Username,
		a.PasswordHash,
		a.IsAdmin,
		a.CreatedAt,
		a.UpdatedAt,
	)
	return account, nil
}

type AccountStorage struct {
	conn postgres.Conn
}

func NewAccountStoragee(conn postgres.Conn) *AccountStorage {
	s := AccountStorage{
		conn: conn,
	}
	return &s
}

func (s *AccountStorage) Create(account *admin.Account) error {
	stmt := `
		INSERT INTO account
			(id, username, password_hash, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5)`

	row, err := marshalAccount(account)
	if err != nil {
		return err
	}

	args := []interface{}{
		row.ID,
		row.Username,
		row.PasswordHash,
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

func (s *AccountStorage) Read(id uuid.UUID) (*admin.Account, error) {
	stmt := `
		SELECT
			id,
			username,
			password_hash,
			is_admin,
			created_at,
			updated_at
		FROM account
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbAccount])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (s *AccountStorage) ReadByUsername(username string) (*admin.Account, error) {
	stmt := `
		SELECT
			id,
			username,
			password_hash,
			is_admin,
			created_at,
			updated_at
		FROM account
		WHERE username = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, username)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbAccount])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (s *AccountStorage) ReadByToken(token string) (*admin.Account, error) {
	stmt := `
		SELECT
			account.id,
			account.username,
			account.password_hash,
			account.is_admin,
			account.created_at,
			account.updated_at
		FROM account
		INNER JOIN token
			ON token.account_id = account.id
		WHERE token.hash = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	hashBytes := sha256.Sum256([]byte(token))
	hash := hex.EncodeToString(hashBytes[:])

	rows, err := s.conn.Query(ctx, stmt, hash)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbAccount])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (repo *AccountStorage) Delete(account *admin.Account) error {
	stmt := `
		DELETE FROM account
		WHERE id = $1
		RETURNING id`

	err := account.CheckDelete()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, account.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}
