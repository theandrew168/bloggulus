package storage_test

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
		account, _ := test.NewAccount(t)
		err := store.Account().Create(account)
		test.AssertNilError(t, err)

		return postgres.ErrRollback
	})
}

func TestAccountCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account, _ := test.CreateAccount(t, store)

		// attempt to create the same account again
		err := store.Account().Create(account)
		test.AssertErrorIs(t, err, postgres.ErrConflict)

		return postgres.ErrRollback
	})
}

func TestAccountRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account, _ := test.CreateAccount(t, store)
		got, err := store.Account().Read(account.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), account.ID())

		return postgres.ErrRollback
	})
}

func TestAccountReadByUsername(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account, _ := test.CreateAccount(t, store)
		got, err := store.Account().ReadByUsername(account.Username())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), account.ID())

		return postgres.ErrRollback
	})
}

func TestAccountReadByToken(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account, _ := test.CreateAccount(t, store)
		_, token := test.CreateToken(t, store, account)

		got, err := store.Account().ReadByToken(token)
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
		account, _ := test.CreateAccount(t, store)

		err := store.Account().Delete(account)
		test.AssertNilError(t, err)

		_, err = store.Account().Read(account.ID())
		test.AssertErrorIs(t, err, postgres.ErrNotFound)

		return postgres.ErrRollback
	})
}
