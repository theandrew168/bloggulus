package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/api/middleware"
	"github.com/theandrew168/bloggulus/backend/web/api/util"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account, _ := test.CreateAccount(t, store)
		_, token := test.CreateToken(t, store, account)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			test.AssertEqual(t, util.ContextGetAccount(r).ID(), account.ID())
		})

		mw := middleware.Authenticate(store)(next)
		mw.ServeHTTP(w, r)

		return postgres.ErrRollback
	})
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

	mw := middleware.Authenticate(store)(next)
	mw.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}

func TestAuthenticateInvalidHeader(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "BearerFOOBAR")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := middleware.Authenticate(store)(next)
	mw.ServeHTTP(w, r)

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

	mw := middleware.Authenticate(store)(next)
	mw.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusUnauthorized)
}
