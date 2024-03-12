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
	logger := test.NewLogger(t)
	storage, closer := test.NewStorage(t)
	defer closer()

	app := api.NewApplication(logger, storage)

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

	for _, test := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", test.url, nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		resp := w.Result()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 400 {
			t.Fatalf("want %v, got %v", 400, resp.StatusCode)
		}

		json := string(body)
		if !strings.Contains(strings.ToLower(json), test.want) {
			t.Fatalf("error JSON missing '%s'", test.want)
		}
	}
}

func TestNotFound(t *testing.T) {
	logger := test.NewLogger(t)
	storage, closer := test.NewStorage(t)
	defer closer()

	app := api.NewApplication(logger, storage)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/missing", nil)

	router := app.Router()
	router.ServeHTTP(w, r)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 404 {
		t.Fatalf("want %v, got %v", 404, resp.StatusCode)
	}

	json := string(body)
	if !strings.Contains(strings.ToLower(json), "not found") {
		t.Fatalf("error JSON missing 'not found'")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	logger := test.NewLogger(t)
	storage, closer := test.NewStorage(t)
	defer closer()

	app := api.NewApplication(logger, storage)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/", nil)

	router := app.Router()
	router.ServeHTTP(w, r)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 405 {
		t.Fatalf("want %v, got %v", 405, resp.StatusCode)
	}

	json := string(body)
	if !strings.Contains(strings.ToLower(json), "method not allowed") {
		t.Fatalf("error JSON missing 'method not allowed'")
	}
}
