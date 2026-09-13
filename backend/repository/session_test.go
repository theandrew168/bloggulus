package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

func TestSessionCreate(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.NewAccount()
	err := repo.Account().Create(context.Background(), account)
	test.AssertNilError(t, err)

	session, _ := test.NewSession(account)
	err = repo.Session().Create(context.Background(), session)
	test.AssertNilError(t, err)
}

func TestSessionCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	session, _ := test.CreateSession(t, repo, account)

	// attempt to create the same session again
	err := repo.Session().Create(context.Background(), session)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestSessionRead(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	session, _ := test.CreateSession(t, repo, account)

	got, err := repo.Session().Read(context.Background(), session.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), session.ID())
}

func TestSessionReadBySessionToken(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	session, sessionToken := test.CreateSession(t, repo, account)

	got, err := repo.Session().ReadByTokenHash(context.Background(), sessionToken.Hash())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), session.ID())
}

func TestSessionDelete(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	session, _ := test.CreateSession(t, repo, account)

	err := repo.Session().Delete(context.Background(), session)
	test.AssertNilError(t, err)

	_, err = repo.Session().Read(context.Background(), session.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}

func TestSessionDeleteExpired(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)

	sessionOld, _, err := model.NewSession(model.NewSessionParams{
		Account: account,
		TTL:     -1 * time.Hour,
	})
	test.AssertNilError(t, err)

	err = repo.Session().Create(context.Background(), sessionOld)
	test.AssertNilError(t, err)

	sessionNew, _, err := model.NewSession(model.NewSessionParams{
		Account: account,
		TTL:     1 * time.Hour,
	})
	test.AssertNilError(t, err)

	err = repo.Session().Create(context.Background(), sessionNew)
	test.AssertNilError(t, err)

	now := timeutil.Now()
	err = repo.Session().DeleteExpired(context.Background(), now)
	test.AssertNilError(t, err)

	_, err = repo.Session().Read(context.Background(), sessionOld.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)

	_, err = repo.Session().Read(context.Background(), sessionNew.ID())
	test.AssertNilError(t, err)
}
