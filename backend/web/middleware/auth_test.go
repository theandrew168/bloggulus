package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	_, sessionID := test.CreateSession(t, repo, account)
	sessionCookie := util.NewSessionCookie(util.SessionCookieName, sessionID)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&sessionCookie)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok := util.GetContextAccount(r)
		test.AssertEqual(t, ok, true)
		test.AssertEqual(t, got.ID(), account.ID())
	})

	h := middleware.Use(next,
		middleware.Authenticate(repo),
	)
	h.ServeHTTP(w, r)
}

func TestAuthenticateNoSession(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := middleware.Use(next, middleware.Authenticate(repo))
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}

func TestAuthenticateInvalidSession(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	sessionCookie := util.NewSessionCookie(util.SessionCookieName, "foobar")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&sessionCookie)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := middleware.Use(next, middleware.Authenticate(repo))
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}

func TestAccountRequired(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	_, sessionID := test.CreateSession(t, repo, account)
	sessionCookie := util.NewSessionCookie(util.SessionCookieName, sessionID)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&sessionCookie)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok := util.GetContextAccount(r)
		test.AssertEqual(t, ok, true)
		test.AssertEqual(t, got.ID(), account.ID())
	})

	h := middleware.Use(next,
		middleware.Authenticate(repo),
		middleware.AccountRequired(),
	)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}

func TestAccountRequiredNoSession(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := middleware.Use(next,
		middleware.Authenticate(repo),
		middleware.AccountRequired(),
	)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusSeeOther)
	test.AssertEqual(t, rr.Header.Get("Location"), "/signin?next=%2F")
}

func TestAccountRequiredInvalidSession(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	sessionCookie := util.NewSessionCookie(util.SessionCookieName, "foobar")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&sessionCookie)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := middleware.Use(next,
		middleware.Authenticate(repo),
		middleware.AccountRequired(),
	)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusSeeOther)
	test.AssertEqual(t, rr.Header.Get("Location"), "/signin?next=%2F")
}

func TestAccountRequiredRedirect(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/foobar", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := middleware.Use(next,
		middleware.Authenticate(repo),
		middleware.AccountRequired(),
	)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusSeeOther)
	test.AssertEqual(t, rr.Header.Get("Location"), "/signin?next=%2Ffoobar")
}

func TestAdminRequired(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	_, sessionID := test.CreateSession(t, repo, account)
	sessionCookie := util.NewSessionCookie(util.SessionCookieName, sessionID)

	// Make the account an admin via manual SQL.
	err := repo.Exec(context.Background(), "UPDATE account SET is_admin = TRUE WHERE id = $1", account.ID())
	test.AssertNilError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&sessionCookie)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok := util.GetContextAccount(r)
		test.AssertEqual(t, ok, true)
		test.AssertEqual(t, got.ID(), account.ID())
	})

	h := middleware.Use(next,
		middleware.Authenticate(repo),
		middleware.AccountRequired(),
		middleware.AdminRequired(),
	)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}

func TestAdminRequiredNoSession(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := middleware.Use(next,
		middleware.Authenticate(repo),
		middleware.AccountRequired(),
		middleware.AdminRequired(),
	)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusSeeOther)
	test.AssertEqual(t, rr.Header.Get("Location"), "/signin?next=%2F")
}

func TestAdminRequiredInvalidSession(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	sessionCookie := util.NewSessionCookie(util.SessionCookieName, "foobar")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&sessionCookie)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := middleware.Use(next,
		middleware.Authenticate(repo),
		middleware.AccountRequired(),
		middleware.AdminRequired(),
	)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusSeeOther)
	test.AssertEqual(t, rr.Header.Get("Location"), "/signin?next=%2F")
}

func TestAdminRequiredNotAdmin(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)
	_, sessionID := test.CreateSession(t, repo, account)
	sessionCookie := util.NewSessionCookie(util.SessionCookieName, sessionID)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&sessionCookie)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := middleware.Use(next,
		middleware.Authenticate(repo),
		middleware.AccountRequired(),
		middleware.AdminRequired(),
	)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusForbidden)
}
