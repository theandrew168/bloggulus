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

func TestAccountRequired(t *testing.T) {
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
			got, ok := util.ContextGetAccount(r)
			test.AssertEqual(t, ok, true)
			test.AssertEqual(t, got.ID(), account.ID())
		})

		mw := middleware.AccountRequired(store)(next)
		mw.ServeHTTP(w, r)

		return postgres.ErrRollback
	})
}

func TestAccountRequiredNoHeader(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := middleware.AccountRequired(store)(next)
	mw.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusUnauthorized)
}

func TestAccountRequiredInvalidHeader(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "BearerFOOBAR")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := middleware.AccountRequired(store)(next)
	mw.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusUnauthorized)
}

func TestAccountRequiredInvalidToken(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer FOOBAR")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := middleware.AccountRequired(store)(next)
	mw.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusUnauthorized)
}

func TestAccountOptional(t *testing.T) {
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
			got, ok := util.ContextGetAccount(r)
			test.AssertEqual(t, ok, true)
			test.AssertEqual(t, got.ID(), account.ID())
		})

		mw := middleware.AccountOptional(store)(next)
		mw.ServeHTTP(w, r)

		return postgres.ErrRollback
	})
}

func TestAccountOptionalNoHeader(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := middleware.AccountOptional(store)(next)
	mw.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}

func TestAccountOptionalInvalidHeader(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "BearerFOOBAR")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := middleware.AccountOptional(store)(next)
	mw.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusUnauthorized)
}

func TestAccountOptionalInvalidToken(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer FOOBAR")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := middleware.AccountOptional(store)(next)
	mw.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusUnauthorized)
}
