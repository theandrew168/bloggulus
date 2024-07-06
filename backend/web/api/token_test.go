package api_test

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
	"github.com/theandrew168/bloggulus/backend/web/api"
)

type jsonNewToken struct {
	ID        uuid.UUID `json:"id"`
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
}

// TODO: For when listing / reading tokens
// type jsonToken struct {
// 	ID        uuid.UUID `json:"id"`
// 	ExpiresAt time.Time `json:"expires_at"`
// }

func TestTokenCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		account, password := test.CreateAccount(t, store)

		h := api.HandleTokenCreate(store)

		req := map[string]string{
			"username": account.Username(),
			"password": password,
		}
		reqBody, err := json.Marshal(req)
		test.AssertNilError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))
		h.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, http.StatusCreated)

		var resp struct {
			Token jsonNewToken `json:"token"`
		}
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		got := resp.Token

		// Ensure the token got created in the database.
		_, err = store.Token().Read(got.ID)
		test.AssertNilError(t, err)

		// Ensure the token can be read by value.
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
		h := api.HandleTokenCreate(store)

		// specify a username that doesn't exist
		req := map[string]string{
			"username": "foo",
			"password": "password",
		}
		reqBody, err := json.Marshal(req)
		test.AssertNilError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))
		h.ServeHTTP(w, r)

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
		account, _ := test.CreateAccount(t, store)

		h := api.HandleTokenCreate(store)

		// specify a password that isn't correct
		req := map[string]string{
			"username": account.Username(),
			"password": "password",
		}
		reqBody, err := json.Marshal(req)
		test.AssertNilError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))
		h.ServeHTTP(w, r)

		rr := w.Result()
		test.AssertEqual(t, rr.StatusCode, http.StatusUnauthorized)

		return postgres.ErrRollback
	})
}
