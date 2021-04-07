package postgres

import (
	"context"

	"github.com/theandrew168/bloggulus/models"
	"github.com/theandrew168/bloggulus/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type sessionStorage struct {
	db *pgxpool.Pool
}

func NewSessionStorage(db *pgxpool.Pool) storage.Session {
	return &sessionStorage{
		db: db,
	}
}

func (s *sessionStorage) Create(ctx context.Context, session *models.Session) (*models.Session, error) {
	command := `
		INSERT INTO session
			(session_id, account_id, expiry)
		VALUES
			($1, $2, $3)`
	_, err := s.db.Exec(ctx, command, session.SessionID, session.AccountID, session.Expiry)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionStorage) Read(ctx context.Context, sessionID string) (*models.Session, error) {
	query := `
		SELECT
			session.*,
			account.*
		FROM session
		INNER JOIN account
			ON account.account_id = session.account_id
		WHERE session_id = $1`
	row := s.db.QueryRow(ctx, query, sessionID)

	var session models.Session
	err := row.Scan(
		&session.SessionID,
		&session.AccountID,
		&session.Expiry,
		&session.Account.AccountID,
		&session.Account.Username,
		&session.Account.Password,
		&session.Account.Email,
		&session.Account.Verified,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *sessionStorage) Delete(ctx context.Context, sessionID string) error {
	command := `
		DELETE
		FROM session
		WHERE session_id = $1`
	_, err := s.db.Exec(ctx, command, sessionID)
	return err
}

func (s *sessionStorage) DeleteExpired(ctx context.Context) error {
	command := `
		DELETE
		FROM session
		WHERE expiry <= now()`
	_ ,err := s.db.Exec(ctx, command)
	return err
}
