package command

import (
	"context"
	"errors"
	"log/slog"
	"uuid"

	"github.com/theandrew168/bloggulus/backend/otel"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
)

var ErrAccountNotFound = errors.New("account: not found")
var ErrDeleteAdminAccount = errors.New("account: cannot delete admin account")

type AccountCommand struct {
	repo *repository.Repository
}

func NewAccount(repo *repository.Repository) *AccountCommand {
	cmd := AccountCommand{
		repo: repo,
	}
	return &cmd
}

func (cmd *AccountCommand) FollowBlog(ctx context.Context, accountID uuid.UUID, blogID uuid.UUID) error {
	ctx, span := otel.GetTracer().Start(ctx, "Command_Account_FollowBlog")
	defer span.End()

	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		account, err := tx.Account().Read(ctx, accountID)
		if err != nil {
			return err
		}

		blog, err := tx.Blog().Read(ctx, blogID)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrBlogNotFound
			}

			return err
		}

		err = account.FollowBlog(blog)
		if err != nil {
			return err
		}

		err = tx.Account().Update(ctx, account)
		if err != nil {
			return err
		}

		slog.InfoContext(ctx, "blog followed",
			"account_id", account.ID().String(),
			"account_username", account.Username().Value(),
			"blog_id", blog.ID().String(),
			"blog_title", blog.Title().Value(),
		)

		return nil
	})
}

func (cmd *AccountCommand) UnfollowBlog(ctx context.Context, accountID uuid.UUID, blogID uuid.UUID) error {
	ctx, span := otel.GetTracer().Start(ctx, "Command_Account_UnfollowBlog")
	defer span.End()

	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		account, err := tx.Account().Read(ctx, accountID)
		if err != nil {
			return err
		}

		blog, err := tx.Blog().Read(ctx, blogID)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrBlogNotFound
			}

			return err
		}

		err = account.UnfollowBlog(blog)
		if err != nil {
			return err
		}

		err = tx.Account().Update(ctx, account)
		if err != nil {
			return err
		}

		slog.InfoContext(ctx, "blog unfollowed",
			"account_id", account.ID().String(),
			"account_username", account.Username().Value(),
			"blog_id", blog.ID().String(),
			"blog_title", blog.Title().Value(),
		)

		return nil
	})
}

func (cmd *AccountCommand) DeleteAccount(ctx context.Context, accountID uuid.UUID) error {
	ctx, span := otel.GetTracer().Start(ctx, "Command_Account_DeleteAccount")
	defer span.End()

	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		account, err := tx.Account().Read(ctx, accountID)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrAccountNotFound
			}

			return err
		}

		// Prevent deletion of admin accounts.
		if account.IsAdmin() {
			return ErrDeleteAdminAccount
		}

		err = tx.Account().Delete(ctx, account)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrAccountNotFound
			}

			return err
		}

		slog.InfoContext(ctx, "account deleted",
			"account_id", account.ID().String(),
			"account_username", account.Username().Value(),
		)

		return nil
	})
}
