package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
)

func TestAddSecureHeaders(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		test.AssertEqual(t, w.Header().Get("Referrer-Policy"), "origin-when-cross-origin")
		test.AssertEqual(t, w.Header().Get("X-Content-Type-Options"), "nosniff")
		test.AssertEqual(t, w.Header().Get("X-Frame-Options"), "deny")
		test.AssertEqual(t, w.Header().Get("X-XSS-Protection"), "0")
	})

	addSecureHeaders := middleware.AddSecureHeaders()
	h := addSecureHeaders(next)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}

func TestEnableCORS(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		test.AssertEqual(t, w.Header().Get("Access-Control-Allow-Origin"), "*")
	})

	enableCORS := middleware.EnableCORS()
	h := enableCORS(next)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}
