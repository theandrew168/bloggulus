package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/internal/api"
	"github.com/theandrew168/bloggulus/internal/core"
	"github.com/theandrew168/bloggulus/internal/postgresql"
	"github.com/theandrew168/bloggulus/internal/test"
)

func TestHandleReadPost(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()

	post := test.CreateMockPost(storage, t)

	app := api.NewApplication(storage, logger)

	url := fmt.Sprintf("/post/%d", post.ID)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", url, nil)

	router := app.Router()
	router.ServeHTTP(w, r)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("want %v, got %v", 200, resp.StatusCode)
	}

	var env map[string]core.Post
	err = json.Unmarshal(body, &env)
	if err != nil {
		t.Fatal(err)
	}

	got, ok := env["post"]
	if !ok {
		t.Fatalf("response missing key: %v", "post")
	}

	if got.ID != post.ID {
		t.Errorf("want %v, got %v", post.ID, got.ID)
	}
}

func TestHandleReadPosts(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()

	test.CreateMockPost(storage, t)

	app := api.NewApplication(storage, logger)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/post", nil)

	router := app.Router()
	router.ServeHTTP(w, r)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("want %v, got %v", 200, resp.StatusCode)
	}

	var env map[string][]core.Post
	err = json.Unmarshal(body, &env)
	if err != nil {
		t.Fatal(err)
	}

	got, ok := env["posts"]
	if !ok {
		t.Fatalf("response missing key: %v", "posts")
	}

	if len(got) < 1 {
		t.Errorf("expected at least one blog")
	}
}

// TODO
func TestHandleReadPostsSearch(t *testing.T) {

}

// TODO
func TestHandleReadPostsPagination(t *testing.T) {

}
