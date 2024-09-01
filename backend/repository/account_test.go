package repository_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestAccountCreate(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account, _ := test.NewAccount(t)
	err := repo.Account().Create(account)
	test.AssertNilError(t, err)
}

func TestAccountCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account, _ := test.CreateAccount(t, repo)

	// attempt to create the same account again
	err := repo.Account().Create(account)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestAccountRead(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account, _ := test.CreateAccount(t, repo)
	got, err := repo.Account().Read(account.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), account.ID())
}

func TestAccountReadByUsername(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account, _ := test.CreateAccount(t, repo)
	got, err := repo.Account().ReadByUsername(account.Username())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), account.ID())
}

func TestAccountReadBySessionIDn(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account, _ := test.CreateAccount(t, repo)
	_, sessionID := test.CreateSession(t, repo, account)

	got, err := repo.Account().ReadBySessionID(sessionID)
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), account.ID())
}

func TestAccountDelete(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account, _ := test.CreateAccount(t, repo)

	err := repo.Account().Delete(account)
	test.AssertNilError(t, err)

	_, err = repo.Account().Read(account.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
