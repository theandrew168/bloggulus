package command

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
)

var ErrAccountNotFound = errors.New("account: not found")
var ErrDeleteAdminAccount = errors.New("account: cannot delete admin account")

func (cmd *Command) FollowBlog(accountID uuid.UUID, blogID uuid.UUID) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		account, err := tx.Account().Read(accountID)
		if err != nil {
			return err
		}

		blog, err := tx.Blog().Read(blogID)
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

		err = tx.Account().Update(account)
		if err != nil {
			return err
		}

		slog.Info("blog followed",
			"account_id", account.ID(),
			"account_username", account.Username(),
			"blog_id", blog.ID(),
			"blog_title", blog.Title(),
		)

		return nil
	})
}

func (cmd *Command) UnfollowBlog(accountID uuid.UUID, blogID uuid.UUID) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		account, err := tx.Account().Read(accountID)
		if err != nil {
			return err
		}

		blog, err := tx.Blog().Read(blogID)
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

		err = tx.Account().Update(account)
		if err != nil {
			return err
		}

		slog.Info("blog unfollowed",
			"account_id", account.ID(),
			"account_username", account.Username(),
			"blog_id", blog.ID(),
			"blog_title", blog.Title(),
		)

		return nil
	})
}

func (cmd *Command) DeleteAccount(accountID uuid.UUID) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		account, err := tx.Account().Read(accountID)
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

		err = tx.Account().Delete(account)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrAccountNotFound
			}

			return err
		}

		return nil
	})
}
