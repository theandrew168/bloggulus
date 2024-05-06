package admin_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/api/admin"
)

type jsonAccount struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func TestAccountCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := admin.NewApplication(store)

		req := map[string]string{
			"username": "foo",
			"password": "password",
		}
		body, err := json.Marshal(req)
		test.AssertNilError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/accounts", bytes.NewReader(body))

		router := app.Router()
		router.ServeHTTP(w, r)

		rr := w.Result()
		body, err = io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, http.StatusCreated)

		var resp map[string]jsonAccount
		err = json.Unmarshal(body, &resp)
		test.AssertNilError(t, err)

		got, ok := resp["account"]
		if !ok {
			t.Fatalf("response missing key: %v", "account")
		}

		test.AssertEqual(t, got.Username, "foo")

		return postgres.ErrRollback
	})
}
