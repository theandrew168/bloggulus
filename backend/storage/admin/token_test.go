package admin_test

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
		account := test.NewAccount(t)
		err := store.Admin().Account().Create(account)
		test.AssertNilError(t, err)

		token := test.NewToken(t, account)
		err = store.Admin().Token().Create(token)
		test.AssertNilError(t, err)

		return postgres.ErrRollback
	})
}

func TestTokenCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		token := test.CreateToken(t, store)

		// attempt to create the same token again
		err := store.Admin().Token().Create(token)
		test.AssertErrorIs(t, err, postgres.ErrConflict)

		return postgres.ErrRollback
	})
}

func TestTokenRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		token := test.CreateToken(t, store)
		got, err := store.Admin().Token().Read(token.ID())
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
		token := test.CreateToken(t, store)

		err := store.Admin().Token().Delete(token)
		test.AssertNilError(t, err)

		_, err = store.Admin().Token().Read(token.ID())
		test.AssertErrorIs(t, err, postgres.ErrNotFound)

		return postgres.ErrRollback
	})
}
