package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

type dbSession struct {
	ID        uuid.UUID `db:"id"`
	AccountID uuid.UUID `db:"account_id"`
	Hash      string    `db:"hash"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func marshalSession(session *model.Session) (dbSession, error) {
	s := dbSession{
		ID:        session.ID(),
		AccountID: session.AccountID(),
		Hash:      session.Hash(),
		ExpiresAt: session.ExpiresAt(),
		CreatedAt: session.CreatedAt(),
		UpdatedAt: session.UpdatedAt(),
	}
	return s, nil
}

func (s dbSession) unmarshal() (*model.Session, error) {
	session := model.LoadSession(
		s.ID,
		s.AccountID,
		s.Hash,
		s.ExpiresAt,
		s.CreatedAt,
		s.UpdatedAt,
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
			(id, account_id, hash, expires_at, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6)`

	row, err := marshalSession(session)
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

	_, err = r.conn.Exec(ctx, stmt, args...)
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
			session.hash,
			session.expires_at,
			session.created_at,
			session.updated_at
		FROM session
		WHERE session.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbSession])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *SessionRepository) ReadBySessionID(sessionID string) (*model.Session, error) {
	stmt := `
		SELECT
			session.id,
			session.account_id,
			session.hash,
			session.expires_at,
			session.created_at,
			session.updated_at
		FROM session
		WHERE session.hash = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	hashBytes := sha256.Sum256([]byte(sessionID))
	hash := hex.EncodeToString(hashBytes[:])

	rows, err := r.conn.Query(ctx, stmt, hash)
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
			session.hash,
			session.expires_at,
			session.created_at,
			session.updated_at
		FROM session
		WHERE session.expires_at <= $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, now)
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

	err := session.CheckDelete()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, session.ID())
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

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	_, err := r.conn.Exec(ctx, stmt, now)
	if err != nil {
		return err
	}

	return nil
}
