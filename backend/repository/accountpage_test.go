package repository_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestAccountPageCreate(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account, _ := test.CreateAccount(t, repo)
	page := test.CreatePage(t, repo)

	err := repo.AccountPage().Create(account, page)
	test.AssertNilError(t, err)
}

func TestAccountPageCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account, _ := test.CreateAccount(t, repo)
	page := test.CreatePage(t, repo)

	err := repo.AccountPage().Create(account, page)
	test.AssertNilError(t, err)

	err = repo.AccountPage().Create(account, page)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestAccountPageDelete(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account, _ := test.CreateAccount(t, repo)
	page := test.CreatePage(t, repo)

	err := repo.AccountPage().Create(account, page)
	test.AssertNilError(t, err)

	err = repo.AccountPage().Delete(account, page)
	test.AssertNilError(t, err)
}

func TestAccountPageDeleteDoesNotExist(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account, _ := test.NewAccount(t)
	page := test.CreatePage(t, repo)

	err := repo.AccountPage().Delete(account, page)
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
