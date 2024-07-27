package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	account, _ := test.CreateAccount(t, store)
	_, token := test.CreateToken(t, store, account)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok := util.ContextGetAccount(r)
		test.AssertEqual(t, ok, true)
		test.AssertEqual(t, got.ID(), account.ID())
	})

	h := middleware.Use(next,
		middleware.Authenticate(store),
	)
	h.ServeHTTP(w, r)
}

func TestAuthenticateNoHeader(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := middleware.Use(next, middleware.Authenticate(store))
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}

func TestAuthenticatelInvalidHeader(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "BearerFOOBAR")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := middleware.Use(next, middleware.Authenticate(store))
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusUnauthorized)
}

func TestAuthenticateInvalidToken(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer FOOBAR")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := middleware.Use(next, middleware.Authenticate(store))
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusUnauthorized)
}

func TestAccountRequired(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	account, _ := test.CreateAccount(t, store)
	_, token := test.CreateToken(t, store, account)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok := util.ContextGetAccount(r)
		test.AssertEqual(t, ok, true)
		test.AssertEqual(t, got.ID(), account.ID())
	})

	h := middleware.Use(next,
		middleware.Authenticate(store),
		middleware.AccountRequired(),
	)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}

func TestAdminRequired(t *testing.T) {
	t.Parallel()

	db, closer := test.NewDatabase(t)
	defer closer()

	store := storage.New(db)

	account, _ := test.CreateAccount(t, store)
	_, token := test.CreateToken(t, store, account)

	// Make the account an admin via manual SQL.
	err := store.Exec(context.Background(), "UPDATE account SET is_admin = TRUE WHERE id = $1", account.ID())
	test.AssertNilError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok := util.ContextGetAccount(r)
		test.AssertEqual(t, ok, true)
		test.AssertEqual(t, got.ID(), account.ID())
	})

	h := middleware.Use(next,
		middleware.Authenticate(store),
		middleware.AccountRequired(),
		middleware.AdminRequired(),
	)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}
