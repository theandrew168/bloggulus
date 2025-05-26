package command

import (
	"errors"
	"log/slog"
	"time"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/timeutil"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

var ErrSessionNotFound = errors.New("session not found")

func (cmd *Command) SignIn(username string) (string, error) {
	// NOTE: Handling state outside the transaciton is the exception, not the rule.
	// This is a special case where a command needs to return a value (the session ID).
	var sessionID string
	err := cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		account, err := tx.Account().ReadByUsername(username)
		if err != nil {
			if !errors.Is(err, postgres.ErrNotFound) {
				return err
			}

			// We need to create a new account at this point.
			account, err = model.NewAccount(username)
			if err != nil {
				return err
			}

			err = tx.Account().Create(account)
			if err != nil {
				return err
			}

			slog.Info("account created",
				"account_id", account.ID(),
			)
		}

		// Create a new session for the account.
		var session *model.Session
		session, sessionID, err = model.NewSession(account, util.SessionCookieTTL)
		if err != nil {
			return err
		}

		err = tx.Session().Create(session)
		if err != nil {
			return err
		}

		slog.Info("account signed in",
			"account_id", account.ID(),
			"session_id", session.ID(),
		)

		return nil
	})

	return sessionID, err
}

func (cmd *Command) SignOut(sessionID string) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		session, err := tx.Session().ReadBySessionID(sessionID)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrSessionNotFound
			}

			return err
		}

		err = tx.Session().Delete(session)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrSessionNotFound
			}

			return err
		}

		return nil
	})
}

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
				// Ignore any "not found" errors here.
				if errors.Is(err, postgres.ErrNotFound) {
					continue
				}

				return err
			}
		}

		return nil
	})
}
