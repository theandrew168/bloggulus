package postgres

import (
	"context"

	"github.com/theandrew168/bloggulus/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type accountBlogStorage struct {
	db *pgxpool.Pool
}

func NewAccountBlogStorage(db *pgxpool.Pool) storage.AccountBlog {
	return &accountBlogStorage{
		db: db,
	}
}

func (s *accountBlogStorage) Follow(ctx context.Context, accountID int, blogID int) error {
	command := `
		INSERT	
		INTO account_blog
			(account_id, blog_id)
		VALUES
			($1, $2)`
	_, err := s.db.Exec(ctx, command, accountID, blogID)
	return err
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
