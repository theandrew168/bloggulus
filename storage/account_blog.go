package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AccountBlog struct {
	db *pgxpool.Pool
}

func NewAccountBlog(db *pgxpool.Pool) *AccountBlog {
	s := AccountBlog{
		db: db,
	}
	return &s
}

func (s *AccountBlog) Follow(ctx context.Context, accountID int, blogID int) error {
	command := `
		INSERT	
		INTO account_blog
			(account_id, blog_id)
		VALUES
			($1, $2)`
	_, err := s.db.Exec(ctx, command, accountID, blogID)
	return err
}

func (s *AccountBlog) Unfollow(ctx context.Context, accountID int, blogID int) error {
	command := `
		DELETE
		FROM account_blog
		WHERE account_id = $1
		AND blog_id = $2`
	_, err := s.db.Exec(ctx, command, accountID, blogID)
	return err
}
