package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/internal/core"
)

type followStorage struct {
	db *pgxpool.Pool
}

func NewFollowStorage(db *pgxpool.Pool) core.FollowStorage {
	s := followStorage{
		db: db,
	}
	return &s
}

func (s *followStorage) Follow(ctx context.Context, accountID int, blogID int) error {
	command := `
		INSERT INTO follow
			(account_id, blog_id)
		VALUES
			($1, $2)`
	_, err := s.db.Exec(ctx, command, accountID, blogID)
	if err != nil {
		// https://github.com/jackc/pgx/wiki/Error-Handling
		// https://github.com/jackc/pgx/issues/474
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return core.ErrExist
			}
		}
		return err
	}

	return nil
}

func (s *followStorage) Unfollow(ctx context.Context, accountID int, blogID int) error {
	command := `
		DELETE FROM follow
		WHERE account_id = $1
		AND blog_id = $2`
	_, err := s.db.Exec(ctx, command, accountID, blogID)
	return err
}
