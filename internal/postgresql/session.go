package postgresql

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/internal/core"
)

type sessionStorage struct {
	conn *pgxpool.Pool
}

func NewSessionStorage(conn *pgxpool.Pool) core.SessionStorage {
	s := sessionStorage{
		conn: conn,
	}
	return &s
}

func (s *sessionStorage) Create(ctx context.Context, session *core.Session) error {
	stmt := `
		INSERT INTO session
			(session_id, expiry, account_id)
		VALUES
			($1, $2, $3)`
	_, err := s.conn.Exec(ctx, stmt,
		session.SessionID,
		session.Expiry,
		session.Account.AccountID)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionStorage) Read(ctx context.Context, sessionID string) (core.Session, error) {
	stmt := `
		SELECT
			session.session_id,
			session.expiry,
			account.account_id,
			account.username,
			account.password,
			account.email,
			account.verified
		FROM session
		INNER JOIN account
			ON account.account_id = session.account_id
		WHERE session_id = $1`
	row := s.conn.QueryRow(ctx, stmt, sessionID)

	var session core.Session
	err := row.Scan(
		&session.SessionID,
		&session.Expiry,
		&session.Account.AccountID,
		&session.Account.Username,
		&session.Account.Password,
		&session.Account.Email,
		&session.Account.Verified,
	)
	if err != nil {
		return core.Session{}, err
	}

	return session, nil
}

func (s *sessionStorage) Delete(ctx context.Context, sessionID string) error {
	stmt := `
		DELETE
		FROM session
		WHERE session_id = $1`
	_, err := s.conn.Exec(ctx, stmt, sessionID)
	return err
}

func (s *sessionStorage) DeleteExpired(ctx context.Context) error {
	stmt := `
		DELETE
		FROM session
		WHERE expiry <= now()`
	_, err := s.conn.Exec(ctx, stmt)
	return err
}
