package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
)

func TestUse(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		test.AssertEqual(t, w.Header().Get("Referrer-Policy"), "origin-when-cross-origin")
		test.AssertEqual(t, w.Header().Get("X-Content-Type-Options"), "nosniff")
		test.AssertEqual(t, w.Header().Get("X-Frame-Options"), "deny")
		test.AssertEqual(t, w.Header().Get("X-XSS-Protection"), "0")
		test.AssertEqual(t, w.Header().Get("Access-Control-Allow-Origin"), "*")
	})

	h := middleware.Use(next,
		middleware.AddSecureHeaders(),
		middleware.EnableCORS(),
	)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}

func TestChain(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		test.AssertEqual(t, w.Header().Get("Referrer-Policy"), "origin-when-cross-origin")
		test.AssertEqual(t, w.Header().Get("X-Content-Type-Options"), "nosniff")
		test.AssertEqual(t, w.Header().Get("X-Frame-Options"), "deny")
		test.AssertEqual(t, w.Header().Get("X-XSS-Protection"), "0")
		test.AssertEqual(t, w.Header().Get("Access-Control-Allow-Origin"), "*")
	})

	chain := middleware.Chain(middleware.AddSecureHeaders(), middleware.EnableCORS())
	h := chain(next)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}
