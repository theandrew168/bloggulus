package command

import (
	"time"

	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

func (cmd *Command) DeleteExpiredSessions(now time.Time) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		now := timeutil.Now()
		expiredSessions, err := tx.Session().ListExpired(now)
		if err != nil {
			return err
		}

		for _, session := range expiredSessions {
			err := tx.Session().Delete(session)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
