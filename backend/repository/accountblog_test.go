package repository_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestAccountBlogCreate(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	blog := test.CreateBlog(t, repo)

	err := repo.AccountBlog().Create(account, blog)
	test.AssertNilError(t, err)
}

func TestAccountBlogCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	blog := test.CreateBlog(t, repo)

	err := repo.AccountBlog().Create(account, blog)
	test.AssertNilError(t, err)

	err = repo.AccountBlog().Create(account, blog)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestAccountBlogDelete(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	blog := test.CreateBlog(t, repo)

	err := repo.AccountBlog().Create(account, blog)
	test.AssertNilError(t, err)

	err = repo.AccountBlog().Delete(account, blog)
	test.AssertNilError(t, err)
}

func TestAccountBlogDeleteDoesNotExist(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.NewAccount(t)
	blog := test.NewBlog(t)

	err := repo.AccountBlog().Delete(account, blog)
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
