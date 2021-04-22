package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/model"
)

type Session struct {
	db *pgxpool.Pool
}

func NewSession(db *pgxpool.Pool) *Session {
	s := Session{
		db: db,
	}
	return &s
}

func (s *Session) Create(ctx context.Context, session *model.Session) (*model.Session, error) {
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

func (s *Session) Read(ctx context.Context, sessionID string) (*model.Session, error) {
	query := `
		SELECT
			session.*,
			account.*
		FROM session
		INNER JOIN account
			ON account.account_id = session.account_id
		WHERE session_id = $1`
	row := s.db.QueryRow(ctx, query, sessionID)

	var session model.Session
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

func (s *Session) Delete(ctx context.Context, sessionID string) error {
	command := `
		DELETE
		FROM session
		WHERE session_id = $1`
	_, err := s.db.Exec(ctx, command, sessionID)
	return err
}

func (s *Session) DeleteExpired(ctx context.Context) error {
	command := `
		DELETE
		FROM session
		WHERE expiry <= now()`
	_ ,err := s.db.Exec(ctx, command)
	return err
}
