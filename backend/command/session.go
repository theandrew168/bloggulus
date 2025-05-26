package command

import (
	"time"

	"github.com/theandrew168/bloggulus/backend/repository"
)

func (cmd *Command) DeleteExpiredSessions(now time.Time) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		return nil
	})
}
