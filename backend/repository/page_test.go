package repository_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestPageCreate(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	page := test.NewPage(t)
	err := repo.Page().Create(page)
	test.AssertNilError(t, err)
}

func TestPageCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	page := test.CreatePage(t, repo)

	// attempt to create the same page again
	err := repo.Page().Create(page)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestPageRead(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	page := test.CreatePage(t, repo)
	got, err := repo.Page().Read(page.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), page.ID())
}

func TestPageReadByURL(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	page := test.CreatePage(t, repo)
	got, err := repo.Page().ReadByURL(page.URL())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), page.ID())
}

func TestPageListByAccount(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	for i := 0; i < 3; i++ {
		page := test.CreatePage(t, repo)
		test.CreateAccountPage(t, repo, account, page)
	}

	limit := 3
	offset := 0
	pages, err := repo.Page().ListByAccount(account, limit, offset)
	test.AssertNilError(t, err)

	test.AssertEqual(t, len(pages), limit)
}

func TestPageCountByAccount(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	for i := 0; i < 3; i++ {
		page := test.CreatePage(t, repo)
		test.CreateAccountPage(t, repo, account, page)
	}

	count, err := repo.Page().CountByAccount(account)
	test.AssertNilError(t, err)

	test.AssertEqual(t, count, 3)
}

func TestPageDelete(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	page := test.CreatePage(t, repo)

	err := repo.Page().Delete(page)
	test.AssertNilError(t, err)

	_, err = repo.Page().Read(page.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
