package repository_test

import (
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

	account := test.NewAccount(t)
	err := repo.Account().Create(account)
	test.AssertNilError(t, err)

	session, _ := test.NewSession(t, account)
	err = repo.Session().Create(session)
	test.AssertNilError(t, err)
}

func TestSessionCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	session, _ := test.CreateSession(t, repo, account)

	// attempt to create the same session again
	err := repo.Session().Create(session)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestSessionRead(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	session, _ := test.CreateSession(t, repo, account)

	got, err := repo.Session().Read(session.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), session.ID())
}

func TestSessionReadBySessionID(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	session, sessionID := test.CreateSession(t, repo, account)

	got, err := repo.Session().ReadBySessionID(sessionID)
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), session.ID())
}

func TestSessionDelete(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	session, _ := test.CreateSession(t, repo, account)

	err := repo.Session().Delete(session)
	test.AssertNilError(t, err)

	_, err = repo.Session().Read(session.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}

func TestSessionDeleteExpired(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)

	sessionOld, _, err := model.NewSession(
		account,
		-1*time.Hour,
	)
	test.AssertNilError(t, err)

	err = repo.Session().Create(sessionOld)
	test.AssertNilError(t, err)

	sessionNew, _, err := model.NewSession(
		account,
		1*time.Hour,
	)
	test.AssertNilError(t, err)

	err = repo.Session().Create(sessionNew)
	test.AssertNilError(t, err)

	now := timeutil.Now()
	err = repo.Session().DeleteExpired(now)
	test.AssertNilError(t, err)

	_, err = repo.Session().Read(sessionOld.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)

	_, err = repo.Session().Read(sessionNew.ID())
	test.AssertNilError(t, err)
}
