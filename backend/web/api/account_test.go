package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

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

	h := api.HandleAccountCreate(store)

	username := test.RandomString(20)
	req := map[string]string{
		"username": username,
		"password": "password",
	}
	reqBody, err := json.Marshal(req)
	test.AssertNilError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/accounts", bytes.NewReader(reqBody))
	h.ServeHTTP(w, r)

	rr := w.Result()
	respBody, err := io.ReadAll(rr.Body)
	test.AssertNilError(t, err)

	test.AssertEqual(t, rr.StatusCode, http.StatusCreated)

	var resp struct {
		Account jsonAccount `json:"account"`
	}
	err = json.Unmarshal(respBody, &resp)
	test.AssertNilError(t, err)

	got := resp.Account
	test.AssertEqual(t, got.Username, username)

	// Ensure the account got created in the database.
	_, err = store.Account().Read(got.ID)
	test.AssertNilError(t, err)
}

func TestAccountCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	h := api.HandleAccountCreate(store)

	username := test.RandomString(20)
	req := map[string]string{
		"username": username,
		"password": "password",
	}
	reqBody, err := json.Marshal(req)
	test.AssertNilError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/accounts", bytes.NewReader(reqBody))
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusCreated)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/accounts", bytes.NewReader(reqBody))
	h.ServeHTTP(w, r)

	rr = w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusUnprocessableEntity)
}
