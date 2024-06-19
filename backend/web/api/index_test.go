package api_test

import (
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/api"
)

func TestHandleIndex(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	app := api.NewApplication(store)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	router := app.Handler()
	router.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, 200)
}
