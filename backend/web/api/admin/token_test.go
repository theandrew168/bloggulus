package admin_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/api/admin"
)

type jsonNewToken struct {
	ID        uuid.UUID `json:"id"`
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
}

type jsonToken struct {
	ID        uuid.UUID `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func TestTokenCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := admin.NewApplication(store)

		account, password := test.CreateAccount(t, store)

		req := map[string]string{
			"username": account.Username(),
			"password": password,
		}
		reqBody, err := json.Marshal(req)
		test.AssertNilError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/tokens", bytes.NewReader(reqBody))

		router := app.Handler()
		router.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, http.StatusCreated)

		var resp map[string]jsonNewToken
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		got, ok := resp["token"]
		if !ok {
			t.Fatalf("response missing key: %v", "token")
		}

		_, err = store.Token().ReadByValue(got.Value)
		test.AssertNilError(t, err)

		return postgres.ErrRollback
	})
}

func TestTokenCreateInvalidUsername(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := admin.NewApplication(store)

		// specify a username that doesn't exist
		req := map[string]string{
			"username": "foo",
			"password": "password",
		}
		reqBody, err := json.Marshal(req)
		test.AssertNilError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/tokens", bytes.NewReader(reqBody))

		router := app.Handler()
		router.ServeHTTP(w, r)

		rr := w.Result()
		test.AssertEqual(t, rr.StatusCode, http.StatusUnauthorized)

		return postgres.ErrRollback
	})
}

func TestTokenCreateInvalidPassword(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := admin.NewApplication(store)

		account, _ := test.CreateAccount(t, store)

		// specify a password that isn't correct
		req := map[string]string{
			"username": account.Username(),
			"password": "password",
		}
		reqBody, err := json.Marshal(req)
		test.AssertNilError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/tokens", bytes.NewReader(reqBody))

		router := app.Handler()
		router.ServeHTTP(w, r)

		rr := w.Result()
		test.AssertEqual(t, rr.StatusCode, http.StatusUnauthorized)

		return postgres.ErrRollback
	})
}
