package repository

import (
	"context"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/value"
)

type dbSession struct {
	ID            uuid.UUID `db:"id"`
	AccountID     uuid.UUID `db:"account_id"`
	TokenHash     string    `db:"token_hash"`
	ExpiresAt     time.Time `db:"expires_at"`
	MetaCreatedAt time.Time `db:"meta_created_at"`
	MetaUpdatedAt time.Time `db:"meta_updated_at"`
}

func marshalSession(session *model.Session) (dbSession, error) {
	s := dbSession{
		ID:            session.ID(),
		AccountID:     session.AccountID(),
		TokenHash:     session.TokenHash().Value(),
		ExpiresAt:     session.ExpiresAt(),
		MetaCreatedAt: session.Meta().CreatedAt(),
		MetaUpdatedAt: session.Meta().UpdatedAt(),
	}
	return s, nil
}

func (s dbSession) unmarshal() (*model.Session, error) {
	tokenHash, err := value.NewTokenHash(s.TokenHash)
	if err != nil {
		return nil, err
	}

	session := model.LoadSession(
		s.ID,
		s.AccountID,
		tokenHash,
		s.ExpiresAt,
		model.LoadMeta(s.MetaCreatedAt, s.MetaUpdatedAt),
	)
	return session, nil
}

type SessionRepository struct {
	conn postgres.Conn
}

func NewSessionRepository(conn postgres.Conn) *SessionRepository {
	r := SessionRepository{
		conn: conn,
	}
	return &r
}

func (r *SessionRepository) Create(session *model.Session) error {
	stmt := `
		INSERT INTO session
			(id, account_id, token_hash, expires_at, meta_created_at, meta_updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6)`

	row, err := marshalSession(session)
	if err != nil {
		return err
	}

	args := []any{
		row.ID,
		row.AccountID,
		row.TokenHash,
		row.ExpiresAt,
		row.MetaCreatedAt,
		row.MetaUpdatedAt,
	}

	_, err = r.conn.Exec(context.Background(), stmt, args...)
	if err != nil {
		return postgres.CheckCreateError(err)
	}

	return nil
}

func (r *SessionRepository) Read(id uuid.UUID) (*model.Session, error) {
	stmt := `
		SELECT
			session.id,
			session.account_id,
			session.token_hash,
			session.expires_at,
			session.meta_created_at,
			session.meta_updated_at
		FROM session
		WHERE session.id = $1`

	rows, err := r.conn.Query(context.Background(), stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbSession])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *SessionRepository) ReadBySessionToken(token value.Token) (*model.Session, error) {
	stmt := `
		SELECT
			session.id,
			session.account_id,
			session.token_hash,
			session.expires_at,
			session.meta_created_at,
			session.meta_updated_at
		FROM session
		WHERE session.token_hash = $1`

	rows, err := r.conn.Query(context.Background(), stmt, token.Hash().Value())
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbSession])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *SessionRepository) ListExpired(now time.Time) ([]*model.Session, error) {
	stmt := `
		SELECT
			session.id,
			session.account_id,
			session.token_hash,
			session.expires_at,
			session.meta_created_at,
			session.meta_updated_at
		FROM session
		WHERE session.expires_at <= $1`

	rows, err := r.conn.Query(context.Background(), stmt, now)
	if err != nil {
		return nil, err
	}

	sessionRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbSession])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	var sessions []*model.Session
	for _, row := range sessionRows {
		session, err := row.unmarshal()
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *SessionRepository) Delete(session *model.Session) error {
	stmt := `
		DELETE FROM session
		WHERE id = $1
		RETURNING id`

	rows, err := r.conn.Query(context.Background(), stmt, session.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}

func (r *SessionRepository) DeleteExpired(now time.Time) error {
	stmt := `
		DELETE FROM session
		WHERE expires_at <= $1`

	_, err := r.conn.Exec(context.Background(), stmt, now)
	if err != nil {
		return err
	}

	return nil
}
