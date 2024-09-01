package query_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

func TestListArticles(t *testing.T) {
	t.Parallel()

	s, sCloser := test.NewStorage(t)
	defer sCloser()

	q, qCloser := test.NewQuery(t)
	defer qCloser()

	blog := test.CreateBlog(t, s)
	test.CreatePost(t, s, blog)

	articles, err := q.ListArticles(1, 0)
	test.AssertNilError(t, err)

	test.AssertEqual(t, len(articles), 1)
}

func TestListArticlesByAccount(t *testing.T) {
	t.Parallel()

	s, sCloser := test.NewStorage(t)
	defer sCloser()

	q, qCloser := test.NewQuery(t)
	defer qCloser()

	followedBlog := test.CreateBlog(t, s)
	test.CreatePost(t, s, followedBlog)
	test.CreatePost(t, s, followedBlog)
	test.CreatePost(t, s, followedBlog)

	unfollowedBlog := test.CreateBlog(t, s)
	test.CreatePost(t, s, unfollowedBlog)
	test.CreatePost(t, s, unfollowedBlog)
	test.CreatePost(t, s, unfollowedBlog)

	account, _ := test.CreateAccount(t, s)
	test.CreateAccountBlog(t, s, account, followedBlog)

	// List posts from blogs followed by this account.
	articles, err := q.ListArticlesByAccount(account, 5, 0)
	test.AssertNilError(t, err)

	// We should only get the three posts associated with the followed blog.
	test.AssertEqual(t, len(articles), 3)
}

func TestSearchArticles(t *testing.T) {
	t.Parallel()

	s, sCloser := test.NewStorage(t)
	defer sCloser()

	q, qCloser := test.NewQuery(t)
	defer qCloser()

	blog := test.NewBlog(t)
	err := s.Blog().Create(blog)
	test.AssertNilError(t, err)

	// create a post about python
	pythonPost, err := model.NewPost(
		blog,
		test.RandomURL(20),
		"Python",
		"content about python",
		timeutil.Now(),
	)
	test.AssertNilError(t, err)

	err = s.Post().Create(pythonPost)
	test.AssertNilError(t, err)

	// create a post about python
	boringPost, err := model.NewPost(
		blog,
		test.RandomURL(20),
		"Boring",
		"content about nothing",
		timeutil.Now(),
	)
	test.AssertNilError(t, err)

	err = s.Post().Create(boringPost)
	test.AssertNilError(t, err)

	// list articles that relate to python
	articles, err := q.SearchArticles("python", 1, 0)
	test.AssertNilError(t, err)

	// should find at least one
	test.AssertEqual(t, len(articles), 1)
}

func TestSearchArticlesByAccount(t *testing.T) {
	t.Parallel()

	s, sCloser := test.NewStorage(t)
	defer sCloser()

	q, qCloser := test.NewQuery(t)
	defer qCloser()

	// Create some followed posts about python.
	followedBlog := test.CreateBlog(t, s)
	for i := 0; i < 3; i++ {
		post, err := model.NewPost(
			followedBlog,
			test.RandomURL(20),
			"Python",
			"content about python",
			timeutil.Now(),
		)
		test.AssertNilError(t, err)

		err = s.Post().Create(post)
		test.AssertNilError(t, err)
	}

	// Create some unfollowed posts about python.
	unfollowedBlog := test.CreateBlog(t, s)
	for i := 0; i < 3; i++ {
		post, err := model.NewPost(
			unfollowedBlog,
			test.RandomURL(20),
			"Python",
			"content about python",
			timeutil.Now(),
		)
		test.AssertNilError(t, err)

		err = s.Post().Create(post)
		test.AssertNilError(t, err)
	}

	account, _ := test.CreateAccount(t, s)
	test.CreateAccountBlog(t, s, account, followedBlog)

	// List posts (from followed blogs) that relate to python.
	articles, err := q.SearchArticlesByAccount(account, "python", 5, 0)
	test.AssertNilError(t, err)

	// Should only return the three posts from followed blogs.
	test.AssertEqual(t, len(articles), 3)
}

func TestCountArticles(t *testing.T) {
	t.Parallel()

	s, sCloser := test.NewStorage(t)
	defer sCloser()

	q, qCloser := test.NewQuery(t)
	defer qCloser()

	blog := test.CreateBlog(t, s)
	test.CreatePost(t, s, blog)
	test.CreatePost(t, s, blog)
	test.CreatePost(t, s, blog)

	count, err := q.CountArticles()
	test.AssertNilError(t, err)

	test.AssertAtLeast(t, count, 3)
}

func TestCountArticlesByAccount(t *testing.T) {
	t.Parallel()

	s, sCloser := test.NewStorage(t)
	defer sCloser()

	q, qCloser := test.NewQuery(t)
	defer qCloser()

	followedBlog := test.CreateBlog(t, s)
	test.CreatePost(t, s, followedBlog)
	test.CreatePost(t, s, followedBlog)
	test.CreatePost(t, s, followedBlog)

	unfollowedBlog := test.CreateBlog(t, s)
	test.CreatePost(t, s, unfollowedBlog)
	test.CreatePost(t, s, unfollowedBlog)
	test.CreatePost(t, s, unfollowedBlog)

	account, _ := test.CreateAccount(t, s)
	test.CreateAccountBlog(t, s, account, followedBlog)

	// We should only count the three posts associated with the followed blog.
	count, err := q.CountArticlesByAccount(account)
	test.AssertNilError(t, err)
	test.AssertEqual(t, count, 3)
}

func TestCountSearchArticles(t *testing.T) {
	t.Parallel()

	s, sCloser := test.NewStorage(t)
	defer sCloser()

	q, qCloser := test.NewQuery(t)
	defer qCloser()

	blog := test.CreateBlog(t, s)

	// create a post about python
	pythonPost, err := model.NewPost(
		blog,
		test.RandomURL(20),
		"Python",
		"content about python",
		timeutil.Now(),
	)
	test.AssertNilError(t, err)

	err = s.Post().Create(pythonPost)
	test.AssertNilError(t, err)

	// create a post about python
	boringPost, err := model.NewPost(
		blog,
		test.RandomURL(20),
		"Boring",
		"content about nothing",
		timeutil.Now(),
	)
	test.AssertNilError(t, err)

	err = s.Post().Create(boringPost)
	test.AssertNilError(t, err)

	// count posts that relate to python
	count, err := q.CountSearchArticles("python")
	test.AssertNilError(t, err)

	// should find at least one
	test.AssertAtLeast(t, count, 1)
}

func TestCountSearchArticlesByAccount(t *testing.T) {
	t.Parallel()

	s, sCloser := test.NewStorage(t)
	defer sCloser()

	q, qCloser := test.NewQuery(t)
	defer qCloser()

	// Create some followed posts about python.
	followedBlog := test.CreateBlog(t, s)
	for i := 0; i < 3; i++ {
		post, err := model.NewPost(
			followedBlog,
			test.RandomURL(20),
			"Python",
			"content about python",
			timeutil.Now(),
		)
		test.AssertNilError(t, err)

		err = s.Post().Create(post)
		test.AssertNilError(t, err)
	}

	// Create some unfollowed posts about python.
	unfollowedBlog := test.CreateBlog(t, s)
	for i := 0; i < 3; i++ {
		post, err := model.NewPost(
			unfollowedBlog,
			test.RandomURL(20),
			"Python",
			"content about python",
			timeutil.Now(),
		)
		test.AssertNilError(t, err)

		err = s.Post().Create(post)
		test.AssertNilError(t, err)
	}

	account, _ := test.CreateAccount(t, s)
	test.CreateAccountBlog(t, s, account, followedBlog)

	// Count posts (from followed blogs) that relate to python.
	count, err := q.CountSearchArticlesByAccount(account, "python")
	test.AssertNilError(t, err)

	// Should only return the three posts from followed blogs.
	test.AssertEqual(t, count, 3)
}
