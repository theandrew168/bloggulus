package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/model"
)

type accountBlogStorage struct {
	db *pgxpool.Pool
}

func NewAccountBlogStorage(db *pgxpool.Pool) model.AccountBlogStorage {
	s := accountBlogStorage{
		db: db,
	}
	return &s
}

func (s *accountBlogStorage) Follow(ctx context.Context, accountID int, blogID int) error {
	command := `
		INSERT	
		INTO account_blog
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
				return model.ErrExist
			}
		}
		return err
	}

	return nil
}

func (s *accountBlogStorage) Unfollow(ctx context.Context, accountID int, blogID int) error {
	command := `
		DELETE
		FROM account_blog
		WHERE account_id = $1
		AND blog_id = $2`
	_, err := s.db.Exec(ctx, command, accountID, blogID)
	return err
}
