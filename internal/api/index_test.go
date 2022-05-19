package api_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/internal/api"
	"github.com/theandrew168/bloggulus/internal/test"
)

func TestHandleIndex(t *testing.T) {
	logger := test.NewLogger(t)
	storage, closer := test.NewStorage(t)
	defer closer()

	app := api.NewApplication(logger, storage)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

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
}
