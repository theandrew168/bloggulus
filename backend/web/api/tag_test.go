package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/api"
)

type jsonTag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func TestHandleTagCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	h := api.HandleTagCreate(store)

	name := test.RandomString(20)
	req := map[string]string{
		"name": name,
	}

	reqBody, err := json.Marshal(req)
	test.AssertNilError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/tags", bytes.NewReader(reqBody))
	h.ServeHTTP(w, r)

	rr := w.Result()
	respBody, err := io.ReadAll(rr.Body)
	test.AssertNilError(t, err)

	test.AssertEqual(t, rr.StatusCode, 200)

	var resp struct {
		Tag jsonTag `json:"tag"`
	}
	err = json.Unmarshal(respBody, &resp)
	test.AssertNilError(t, err)

	got := resp.Tag
	test.AssertEqual(t, got.Name, name)

	// Ensure the tag got created in the database.
	_, err = store.Tag().Read(got.ID)
	test.AssertNilError(t, err)
}

func TestHandleTagList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	test.CreateTag(t, store)

	h := api.HandleTagList(store)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/tags", nil)
	h.ServeHTTP(w, r)

	rr := w.Result()
	respBody, err := io.ReadAll(rr.Body)
	test.AssertNilError(t, err)

	test.AssertEqual(t, rr.StatusCode, 200)

	var resp struct {
		Count int       `json:"count"`
		Tags  []jsonTag `json:"tags"`
	}
	err = json.Unmarshal(respBody, &resp)
	test.AssertNilError(t, err)

	test.AssertAtLeast(t, resp.Count, 1)
	test.AssertAtLeast(t, len(resp.Tags), 1)
}

func TestHandleTagListPagination(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	// create 5 tags to test with
	test.CreateTag(t, store)
	test.CreateTag(t, store)
	test.CreateTag(t, store)
	test.CreateTag(t, store)
	test.CreateTag(t, store)

	tests := []struct {
		size int
		want int
	}{
		{1, 1},
		{3, 3},
		{5, 5},
	}

	h := api.HandleTagList(store)

	for _, tt := range tests {
		url := fmt.Sprintf("/tags?size=%d", tt.size)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", url, nil)
		h.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp struct {
			Tags []jsonTag `json:"tags"`
		}
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		got := resp.Tags
		test.AssertEqual(t, len(got), tt.want)
	}
}

func TestHandleTagDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	tag := test.CreateTag(t, store)

	h := api.HandleTagDelete(store)

	url := fmt.Sprintf("/tags/%s", tag.ID())
	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", url, nil)
	r.SetPathValue("tagID", tag.ID().String())
	h.ServeHTTP(w, r)

	rr := w.Result()
	respBody, err := io.ReadAll(rr.Body)
	test.AssertNilError(t, err)

	test.AssertEqual(t, rr.StatusCode, 200)

	var resp struct {
		Tag jsonTag `json:"tag"`
	}
	err = json.Unmarshal(respBody, &resp)
	test.AssertNilError(t, err)

	got := resp.Tag
	test.AssertEqual(t, got.ID, tag.ID())

	_, err = store.Tag().Read(got.ID)
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
