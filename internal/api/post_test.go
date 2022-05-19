package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/theandrew168/bloggulus"
	"github.com/theandrew168/bloggulus/internal/api"
	"github.com/theandrew168/bloggulus/internal/test"
)

func TestHandleReadPost(t *testing.T) {
	logger := test.NewLogger(t)
	storage, closer := test.NewStorage(t)
	defer closer()

	app := api.NewApplication(logger, storage)

	post := test.CreateMockPost(t, storage)

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
		t.Fatalf("want %v, got %v", 200, resp.StatusCode)
	}

	var env map[string]bloggulus.Post
	err = json.Unmarshal(body, &env)
	if err != nil {
		t.Fatal(err)
	}

	got, ok := env["post"]
	if !ok {
		t.Fatalf("response missing key: %v", "post")
	}

	if got.ID != post.ID {
		t.Fatalf("want %v, got %v", post.ID, got.ID)
	}
}

func TestHandleReadPostNotFound(t *testing.T) {
	logger := test.NewLogger(t)
	storage, closer := test.NewStorage(t)
	defer closer()

	app := api.NewApplication(logger, storage)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/post/999999999", nil)

	router := app.Router()
	router.ServeHTTP(w, r)

	resp := w.Result()
	if resp.StatusCode != 404 {
		t.Fatalf("want %v, got %v", 404, resp.StatusCode)
	}
}

func TestHandleReadPosts(t *testing.T) {
	logger := test.NewLogger(t)
	storage, closer := test.NewStorage(t)
	defer closer()

	app := api.NewApplication(logger, storage)

	test.CreateMockPost(t, storage)

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
		t.Fatalf("want %v, got %v", 200, resp.StatusCode)
	}

	var env map[string][]bloggulus.Post
	err = json.Unmarshal(body, &env)
	if err != nil {
		t.Fatal(err)
	}

	got, ok := env["posts"]
	if !ok {
		t.Fatalf("response missing key: %v", "posts")
	}

	if len(got) < 1 {
		t.Fatalf("expected at least one blog")
	}
}

func TestHandleReadPostsPagination(t *testing.T) {
	logger := test.NewLogger(t)
	storage, closer := test.NewStorage(t)
	defer closer()

	app := api.NewApplication(logger, storage)

	// create 5 posts to test with
	test.CreateMockPost(t, storage)
	test.CreateMockPost(t, storage)
	test.CreateMockPost(t, storage)
	test.CreateMockPost(t, storage)
	test.CreateMockPost(t, storage)

	tests := []struct {
		limit int
		want  int
	}{
		{0, 0},
		{1, 1},
		{3, 3},
		{5, 5},
	}

	for _, test := range tests {
		url := fmt.Sprintf("/post?limit=%d", test.limit)
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
			t.Fatalf("want %v, got %v", 200, resp.StatusCode)
		}

		var env map[string][]bloggulus.Post
		err = json.Unmarshal(body, &env)
		if err != nil {
			t.Fatal(err)
		}

		got, ok := env["posts"]
		if !ok {
			t.Fatalf("response missing key: %v", "posts")
		}

		if len(got) != test.want {
			t.Fatalf("want %v, got %v", test.want, len(got))
		}
	}
}

func TestHandleReadPostsSearch(t *testing.T) {
	logger := test.NewLogger(t)
	storage, closer := test.NewStorage(t)
	defer closer()

	app := api.NewApplication(logger, storage)

	blog := test.CreateMockBlog(t, storage)
	q := "python rust"

	// create searchable post
	post := bloggulus.NewPost(test.RandomURL(32), q, test.RandomTime(), blog)
	err := storage.Post.Create(&post)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("/post?q=%s", url.QueryEscape(q))
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
		t.Fatalf("want %v, got %v", 200, resp.StatusCode)
	}

	var env map[string][]bloggulus.Post
	err = json.Unmarshal(body, &env)
	if err != nil {
		t.Fatal(err)
	}

	got, ok := env["posts"]
	if !ok {
		t.Fatalf("response missing key: %v", "posts")
	}

	if len(got) < 1 {
		t.Fatalf("expected at least one matching post")
	}
}
