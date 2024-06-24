package storage_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestTokenCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account, _ := test.NewAccount(t)
		err := store.Account().Create(account)
		test.AssertNilError(t, err)

		token, _ := test.NewToken(t, account)
		err = store.Token().Create(token)
		test.AssertNilError(t, err)

		return postgres.ErrRollback
	})
}

func TestTokenCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account, _ := test.CreateAccount(t, store)
		token, _ := test.CreateToken(t, store, account)

		// attempt to create the same token again
		err := store.Token().Create(token)
		test.AssertErrorIs(t, err, postgres.ErrConflict)

		return postgres.ErrRollback
	})
}

func TestTokenRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account, _ := test.CreateAccount(t, store)
		token, _ := test.CreateToken(t, store, account)

		got, err := store.Token().Read(token.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), token.ID())

		return postgres.ErrRollback
	})
}

func TestTokenReadByValue(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account, _ := test.CreateAccount(t, store)
		token, value := test.CreateToken(t, store, account)

		got, err := store.Token().ReadByValue(value)
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), token.ID())

		return postgres.ErrRollback
	})
}

func TestTokenDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account, _ := test.CreateAccount(t, store)
		token, _ := test.CreateToken(t, store, account)

		err := store.Token().Delete(token)
		test.AssertNilError(t, err)

		_, err = store.Token().Read(token.ID())
		test.AssertErrorIs(t, err, postgres.ErrNotFound)

		return postgres.ErrRollback
	})
}
