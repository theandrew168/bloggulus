package web_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/internal/postgresql"
	"github.com/theandrew168/bloggulus/internal/test"
	"github.com/theandrew168/bloggulus/internal/web"
)

func TestHandleIndex(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()

	app := web.NewApplication(storage, logger)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	// TODO: test the handler independent of its route?
	router := app.Router()
	router.ServeHTTP(w, r)

	resp := w.Result()
	_, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("want %v, got %v\n", 200, resp.StatusCode)
	}

	// TODO: test particular elements within body?
	// TODO: how to do this without writing flaky tests that depend on design?
}
