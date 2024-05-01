package admin_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestAccountCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account := test.NewAccount(t)
		err := store.Admin().Account().Create(account)
		test.AssertNilError(t, err)

		return postgres.ErrRollback
	})
}

func TestAccountCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account := test.CreateAccount(t, store)

		// attempt to create the same account again
		err := store.Admin().Account().Create(account)
		test.AssertErrorIs(t, err, postgres.ErrConflict)

		return postgres.ErrRollback
	})
}

func TestAccountRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account := test.CreateAccount(t, store)
		got, err := store.Admin().Account().Read(account.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), account.ID())

		return postgres.ErrRollback
	})
}

func TestAccountDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account := test.CreateAccount(t, store)

		err := store.Admin().Account().Delete(account)
		test.AssertNilError(t, err)

		_, err = store.Admin().Account().Read(account.ID())
		test.AssertErrorIs(t, err, postgres.ErrNotFound)

		return postgres.ErrRollback
	})
}
