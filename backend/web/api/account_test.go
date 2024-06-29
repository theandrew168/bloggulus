package api_test

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
	"github.com/theandrew168/bloggulus/backend/web/api"
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
		syncService := test.NewSyncService(t, store, nil, nil)

		app := api.NewApplication(store, syncService)

		req := map[string]string{
			"username": "foo",
			"password": "password",
		}
		reqBody, err := json.Marshal(req)
		test.AssertNilError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/accounts", bytes.NewReader(reqBody))

		router := app.Handler()
		router.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, http.StatusCreated)

		var resp map[string]jsonAccount
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		got, ok := resp["account"]
		if !ok {
			t.Fatalf("response missing key: %v", "account")
		}

		test.AssertEqual(t, got.Username, "foo")

		// Ensure the account got created in the database.
		_, err = store.Account().Read(got.ID)
		test.AssertNilError(t, err)

		return postgres.ErrRollback
	})
}

func TestAccountCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		syncService := test.NewSyncService(t, store, nil, nil)

		app := api.NewApplication(store, syncService)
		router := app.Handler()

		req := map[string]string{
			"username": "foo",
			"password": "password",
		}
		reqBody, err := json.Marshal(req)
		test.AssertNilError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/accounts", bytes.NewReader(reqBody))
		router.ServeHTTP(w, r)

		rr := w.Result()
		test.AssertEqual(t, rr.StatusCode, http.StatusCreated)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("POST", "/accounts", bytes.NewReader(reqBody))
		router.ServeHTTP(w, r)

		rr = w.Result()
		test.AssertEqual(t, rr.StatusCode, http.StatusUnprocessableEntity)

		return postgres.ErrRollback
	})
}
