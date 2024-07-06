package api_test

import (
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/api"
)

func TestHandleIndexRedoc(t *testing.T) {
	t.Parallel()

	h := api.HandleIndexRedoc()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/redoc", nil)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, 200)
}

func TestHandleIndexRapidoc(t *testing.T) {
	t.Parallel()

	h := api.HandleIndexRapidoc()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/rapidoc", nil)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, 200)
}

func TestHandleIndexStoplight(t *testing.T) {
	t.Parallel()

	h := api.HandleIndexStoplight()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/stoplight", nil)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, 200)
}
