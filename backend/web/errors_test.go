package web_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web"
)

func TestNotFound(t *testing.T) {
	logger := test.NewLogger(t)
	storage, closer := test.NewStorage(t)
	defer closer()

	app := web.NewApplication(logger, storage)

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

	page := string(body)
	if !strings.Contains(strings.ToLower(page), "not found") {
		t.Fatalf("error page missing 'not found'")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	logger := test.NewLogger(t)
	storage, closer := test.NewStorage(t)
	defer closer()

	app := web.NewApplication(logger, storage)

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

	page := string(body)
	if !strings.Contains(strings.ToLower(page), "method not allowed") {
		t.Fatalf("error page missing 'method not allowed'")
	}
}
