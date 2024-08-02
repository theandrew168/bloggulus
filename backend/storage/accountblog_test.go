package storage_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestAccountBlogCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	account, _ := test.CreateAccount(t, store)
	blog := test.CreateBlog(t, store)

	err := store.AccountBlog().Create(account, blog)
	test.AssertNilError(t, err)
}

func TestAccountBlogCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	account, _ := test.CreateAccount(t, store)
	blog := test.CreateBlog(t, store)

	err := store.AccountBlog().Create(account, blog)
	test.AssertNilError(t, err)

	err = store.AccountBlog().Create(account, blog)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestAccountBlogDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	account, _ := test.CreateAccount(t, store)
	blog := test.CreateBlog(t, store)

	err := store.AccountBlog().Create(account, blog)
	test.AssertNilError(t, err)

	err = store.AccountBlog().Delete(account, blog)
	test.AssertNilError(t, err)
}

func TestAccountBlogDeleteDoesNotExist(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	account, _ := test.NewAccount(t)
	blog := test.NewBlog(t)

	err := store.AccountBlog().Delete(account, blog)
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
