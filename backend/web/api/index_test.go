package api_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/api"
)

func TestHandleIndex(t *testing.T) {
	t.Parallel()

	storage, closer := test.NewAdminStorage(t)
	defer closer()

	app := api.NewApplication(storage)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	router := app.Router()
	router.ServeHTTP(w, r)

	rr := w.Result()
	_, err := io.ReadAll(rr.Body)
	test.AssertNilError(t, err)

	test.AssertEqual(t, rr.StatusCode, 200)
}
