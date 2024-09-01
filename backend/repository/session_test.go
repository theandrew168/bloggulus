package repository_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestSessionCreate(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account, _ := test.NewAccount(t)
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

	account, _ := test.CreateAccount(t, repo)
	session, _ := test.CreateSession(t, repo, account)

	// attempt to create the same session again
	err := repo.Session().Create(session)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestSessionRead(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account, _ := test.CreateAccount(t, repo)
	session, _ := test.CreateSession(t, repo, account)

	got, err := repo.Session().Read(session.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), session.ID())
}

func TestSessionReadBySessionID(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account, _ := test.CreateAccount(t, repo)
	session, sessionID := test.CreateSession(t, repo, account)

	got, err := repo.Session().ReadBySessionID(sessionID)
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), session.ID())
}

func TestSessionDelete(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account, _ := test.CreateAccount(t, repo)
	session, _ := test.CreateSession(t, repo, account)

	err := repo.Session().Delete(session)
	test.AssertNilError(t, err)

	_, err = repo.Session().Read(session.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
