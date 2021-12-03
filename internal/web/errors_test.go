package web_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theandrew168/bloggulus/internal/postgresql"
	"github.com/theandrew168/bloggulus/internal/test"
	"github.com/theandrew168/bloggulus/internal/web"
)

func TestNotFound(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()
	app := web.NewApplication(storage, logger)

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
		t.Errorf("want %v, got %v", 404, resp.StatusCode)
	}

	page := string(body)
	if !strings.Contains(strings.ToLower(page), "not found") {
		t.Errorf("error page missing 'not found'")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()
	app := web.NewApplication(storage, logger)

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
		t.Errorf("want %v, got %v", 405, resp.StatusCode)
	}

	page := string(body)
	if !strings.Contains(strings.ToLower(page), "method not allowed") {
		t.Errorf("error page missing 'method not allowed'")
	}
}
