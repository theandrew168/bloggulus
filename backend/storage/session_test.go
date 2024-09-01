package storage_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestSessionCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	account, _ := test.NewAccount(t)
	err := store.Account().Create(account)
	test.AssertNilError(t, err)

	session, _ := test.NewSession(t, account)
	err = store.Session().Create(session)
	test.AssertNilError(t, err)
}

func TestSessionCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	account, _ := test.CreateAccount(t, store)
	session, _ := test.CreateSession(t, store, account)

	// attempt to create the same session again
	err := store.Session().Create(session)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestSessionRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	account, _ := test.CreateAccount(t, store)
	session, _ := test.CreateSession(t, store, account)

	got, err := store.Session().Read(session.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), session.ID())
}

func TestSessionDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	account, _ := test.CreateAccount(t, store)
	session, _ := test.CreateSession(t, store, account)

	err := store.Session().Delete(session)
	test.AssertNilError(t, err)

	_, err = store.Session().Read(session.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
