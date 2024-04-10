package api_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/api"
)

func TestBadRequest(t *testing.T) {
	t.Parallel()

	storage, closer := test.NewAdminStorage(t)
	defer closer()

	app := api.NewApplication(storage)

	tests := []struct {
		url  string
		want string
	}{
		{"/blogs/invalid", "uuid"},
		{"/blogs?limit=asdf", "integer"},
		{"/blogs?limit=-123", "positive"},
		{"/blogs?limit=123", "less than"},
		{"/blogs?offset=asdf", "integer"},
		{"/blogs?offset=-123", "positive"},
		{"/posts/invalid", "uuid"},
		{"/posts?limit=asdf", "integer"},
		{"/posts?limit=-123", "positive"},
		{"/posts?limit=123", "less than"},
		{"/posts?offset=asdf", "integer"},
		{"/posts?offset=-123", "positive"},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", tt.url, nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		rr := w.Result()
		body, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 400)

		json := string(body)
		test.AssertStringContains(t, strings.ToLower(json), tt.want)
	}
}

func TestNotFound(t *testing.T) {
	t.Parallel()

	storage, closer := test.NewAdminStorage(t)
	defer closer()

	app := api.NewApplication(storage)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/missing", nil)

	router := app.Router()
	router.ServeHTTP(w, r)

	rr := w.Result()
	body, err := io.ReadAll(rr.Body)
	test.AssertNilError(t, err)

	test.AssertEqual(t, rr.StatusCode, 404)

	json := string(body)
	test.AssertStringContains(t, strings.ToLower(json), "not found")
}

func TestMethodNotAllowed(t *testing.T) {
	t.Parallel()

	storage, closer := test.NewAdminStorage(t)
	defer closer()

	app := api.NewApplication(storage)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/", nil)

	router := app.Router()
	router.ServeHTTP(w, r)

	rr := w.Result()
	body, err := io.ReadAll(rr.Body)
	test.AssertNilError(t, err)

	test.AssertEqual(t, rr.StatusCode, 405)

	json := string(body)
	test.AssertStringContains(t, strings.ToLower(json), "method not allowed")
}
