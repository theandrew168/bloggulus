package command

import (
	"errors"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
)

var ErrAccountNotFound = errors.New("account not found")
var ErrDeleteAdminAccount = errors.New("cannot delete admin account")

func (cmd *Command) FollowBlog(accountID uuid.UUID, blogID uuid.UUID) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		account, err := tx.Account().Read(accountID)
		if err != nil {
			return err
		}

		blog, err := tx.Blog().Read(blogID)
		if err != nil {
			return err
		}

		err = account.FollowBlog(blog)
		if err != nil {
			return err
		}

		return tx.Account().Update(account)
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
			return err
		}

		err = account.UnfollowBlog(blog)
		if err != nil {
			return err
		}

		return tx.Account().Update(account)
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
