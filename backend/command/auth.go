package command

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/timeutil"
	"github.com/theandrew168/bloggulus/backend/value"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

var ErrSessionNotFound = errors.New("session: not found")

type AuthCommand struct {
	repo *repository.Repository
}

func NewAuth(repo *repository.Repository) *AuthCommand {
	cmd := AuthCommand{
		repo: repo,
	}
	return &cmd
}

func (cmd *AuthCommand) SignIn(ctx context.Context, username value.Name) (value.Token, error) {
	// NOTE: Handling state outside the transaction is the exception, not the rule.
	// This is a special case where a command needs to return a value (the session ID).
	var sessionToken value.Token
	err := cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		account, err := tx.Account().ReadByUsername(ctx, username)
		if err != nil {
			if !errors.Is(err, postgres.ErrNotFound) {
				return err
			}

			// We need to create a new account at this point.
			account, err = model.NewAccount(model.NewAccountParams{
				Username: username,
			})
			if err != nil {
				return err
			}

			err = tx.Account().Create(ctx, account)
			if err != nil {
				return err
			}

			slog.Info("account created",
				"account_id", account.ID().String(),
			)
		}

		// Create a new session for the account.
		var session *model.Session
		session, sessionToken, err = model.NewSession(model.NewSessionParams{
			Account: account,
			TTL:     util.SessionCookieTTL,
		})
		if err != nil {
			return err
		}

		err = tx.Session().Create(ctx, session)
		if err != nil {
			return err
		}

		slog.Info("account signed in",
			"account_id", account.ID().String(),
			"session_id", session.ID().String(),
		)

		return nil
	})

	return sessionToken, err
}

func (cmd *AuthCommand) SignOut(ctx context.Context, sessionToken value.Token) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		session, err := tx.Session().ReadByTokenHash(ctx, sessionToken.Hash())
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrSessionNotFound
			}

			return err
		}

		err = tx.Session().Delete(ctx, session)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrSessionNotFound
			}

			return err
		}

		return nil
	})
}

func (cmd *AuthCommand) DeleteExpiredSessions(ctx context.Context, now time.Time) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		now := timeutil.Now()
		expiredSessions, err := tx.Session().ListExpired(ctx, now)
		if err != nil {
			return err
		}

		for _, session := range expiredSessions {
			err := tx.Session().Delete(ctx, session)
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
