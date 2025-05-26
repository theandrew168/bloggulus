package command

import (
	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/repository"
)

func (cmd *Command) FollowBlog(accountID uuid.UUID, blogID uuid.UUID) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		account, err := tx.Account().Read(accountID)
		if err != nil {
			return err
		}

		err = account.FollowBlog(blogID)
		if err != nil {
			return err
		}

		return tx.Account().Update(account)
	})
}
