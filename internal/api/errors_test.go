package api_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theandrew168/bloggulus/internal/api"
	"github.com/theandrew168/bloggulus/internal/postgresql"
	"github.com/theandrew168/bloggulus/internal/test"
)

func TestBadRequest(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()
	app := api.NewApplication(storage, logger)

	tests := []struct {
		url  string
		want string
	}{
		{"/blog/invalid", "integer"},
		{"/blog/-123", "positive"},
		{"/blog?limit=asdf", "integer"},
		{"/blog?limit=-123", "positive"},
		{"/blog?limit=123", "less than"},
		{"/blog?offset=asdf", "integer"},
		{"/blog?offset=-123", "positive"},
		{"/post/invalid", "integer"},
		{"/post/-123", "positive"},
		{"/post?limit=asdf", "integer"},
		{"/post?limit=-123", "positive"},
		{"/post?limit=123", "less than"},
		{"/post?offset=asdf", "integer"},
		{"/post?offset=-123", "positive"},
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
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()
	app := api.NewApplication(storage, logger)

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
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()
	app := api.NewApplication(storage, logger)

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
