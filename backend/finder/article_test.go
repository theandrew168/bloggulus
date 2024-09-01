package finder_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

func TestListArticles(t *testing.T) {
	t.Parallel()

	repo, repoCloser := test.NewRepository(t)
	defer repoCloser()

	find, findCloser := test.NewFinder(t)
	defer findCloser()

	blog := test.CreateBlog(t, repo)
	test.CreatePost(t, repo, blog)

	articles, err := find.ListArticles(1, 0)
	test.AssertNilError(t, err)

	test.AssertEqual(t, len(articles), 1)
}

func TestListArticlesByAccount(t *testing.T) {
	t.Parallel()

	repo, repoCloser := test.NewRepository(t)
	defer repoCloser()

	find, findCloser := test.NewFinder(t)
	defer findCloser()

	followedBlog := test.CreateBlog(t, repo)
	test.CreatePost(t, repo, followedBlog)
	test.CreatePost(t, repo, followedBlog)
	test.CreatePost(t, repo, followedBlog)

	unfollowedBlog := test.CreateBlog(t, repo)
	test.CreatePost(t, repo, unfollowedBlog)
	test.CreatePost(t, repo, unfollowedBlog)
	test.CreatePost(t, repo, unfollowedBlog)

	account, _ := test.CreateAccount(t, repo)
	test.CreateAccountBlog(t, repo, account, followedBlog)

	// List posts from blogs followed by this account.
	articles, err := find.ListArticlesByAccount(account, 5, 0)
	test.AssertNilError(t, err)

	// We should only get the three posts associated with the followed blog.
	test.AssertEqual(t, len(articles), 3)
}

func TestSearchArticles(t *testing.T) {
	t.Parallel()

	repo, repoCloser := test.NewRepository(t)
	defer repoCloser()

	find, findCloser := test.NewFinder(t)
	defer findCloser()

	blog := test.NewBlog(t)
	err := repo.Blog().Create(blog)
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

	err = repo.Post().Create(pythonPost)
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

	err = repo.Post().Create(boringPost)
	test.AssertNilError(t, err)

	// list articles that relate to python
	articles, err := find.SearchArticles("python", 1, 0)
	test.AssertNilError(t, err)

	// should find at least one
	test.AssertEqual(t, len(articles), 1)
}

func TestSearchArticlesByAccount(t *testing.T) {
	t.Parallel()

	repo, repoCloser := test.NewRepository(t)
	defer repoCloser()

	find, findCloser := test.NewFinder(t)
	defer findCloser()

	// Create some followed posts about python.
	followedBlog := test.CreateBlog(t, repo)
	for i := 0; i < 3; i++ {
		post, err := model.NewPost(
			followedBlog,
			test.RandomURL(20),
			"Python",
			"content about python",
			timeutil.Now(),
		)
		test.AssertNilError(t, err)

		err = repo.Post().Create(post)
		test.AssertNilError(t, err)
	}

	// Create some unfollowed posts about python.
	unfollowedBlog := test.CreateBlog(t, repo)
	for i := 0; i < 3; i++ {
		post, err := model.NewPost(
			unfollowedBlog,
			test.RandomURL(20),
			"Python",
			"content about python",
			timeutil.Now(),
		)
		test.AssertNilError(t, err)

		err = repo.Post().Create(post)
		test.AssertNilError(t, err)
	}

	account, _ := test.CreateAccount(t, repo)
	test.CreateAccountBlog(t, repo, account, followedBlog)

	// List posts (from followed blogs) that relate to python.
	articles, err := find.SearchArticlesByAccount(account, "python", 5, 0)
	test.AssertNilError(t, err)

	// Should only return the three posts from followed blogs.
	test.AssertEqual(t, len(articles), 3)
}

func TestCountArticles(t *testing.T) {
	t.Parallel()

	repo, repoCloser := test.NewRepository(t)
	defer repoCloser()

	find, findCloser := test.NewFinder(t)
	defer findCloser()

	blog := test.CreateBlog(t, repo)
	test.CreatePost(t, repo, blog)
	test.CreatePost(t, repo, blog)
	test.CreatePost(t, repo, blog)

	count, err := find.CountArticles()
	test.AssertNilError(t, err)

	test.AssertAtLeast(t, count, 3)
}

func TestCountArticlesByAccount(t *testing.T) {
	t.Parallel()

	repo, repoCloser := test.NewRepository(t)
	defer repoCloser()

	find, findCloser := test.NewFinder(t)
	defer findCloser()

	followedBlog := test.CreateBlog(t, repo)
	test.CreatePost(t, repo, followedBlog)
	test.CreatePost(t, repo, followedBlog)
	test.CreatePost(t, repo, followedBlog)

	unfollowedBlog := test.CreateBlog(t, repo)
	test.CreatePost(t, repo, unfollowedBlog)
	test.CreatePost(t, repo, unfollowedBlog)
	test.CreatePost(t, repo, unfollowedBlog)

	account, _ := test.CreateAccount(t, repo)
	test.CreateAccountBlog(t, repo, account, followedBlog)

	// We should only count the three posts associated with the followed blog.
	count, err := find.CountArticlesByAccount(account)
	test.AssertNilError(t, err)
	test.AssertEqual(t, count, 3)
}

func TestCountSearchArticles(t *testing.T) {
	t.Parallel()

	repo, repoCloser := test.NewRepository(t)
	defer repoCloser()

	find, findCloser := test.NewFinder(t)
	defer findCloser()

	blog := test.CreateBlog(t, repo)

	// create a post about python
	pythonPost, err := model.NewPost(
		blog,
		test.RandomURL(20),
		"Python",
		"content about python",
		timeutil.Now(),
	)
	test.AssertNilError(t, err)

	err = repo.Post().Create(pythonPost)
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

	err = repo.Post().Create(boringPost)
	test.AssertNilError(t, err)

	// count posts that relate to python
	count, err := find.CountSearchArticles("python")
	test.AssertNilError(t, err)

	// should find at least one
	test.AssertAtLeast(t, count, 1)
}

func TestCountSearchArticlesByAccount(t *testing.T) {
	t.Parallel()

	repo, repoCloser := test.NewRepository(t)
	defer repoCloser()

	find, findCloser := test.NewFinder(t)
	defer findCloser()

	// Create some followed posts about python.
	followedBlog := test.CreateBlog(t, repo)
	for i := 0; i < 3; i++ {
		post, err := model.NewPost(
			followedBlog,
			test.RandomURL(20),
			"Python",
			"content about python",
			timeutil.Now(),
		)
		test.AssertNilError(t, err)

		err = repo.Post().Create(post)
		test.AssertNilError(t, err)
	}

	// Create some unfollowed posts about python.
	unfollowedBlog := test.CreateBlog(t, repo)
	for i := 0; i < 3; i++ {
		post, err := model.NewPost(
			unfollowedBlog,
			test.RandomURL(20),
			"Python",
			"content about python",
			timeutil.Now(),
		)
		test.AssertNilError(t, err)

		err = repo.Post().Create(post)
		test.AssertNilError(t, err)
	}

	account, _ := test.CreateAccount(t, repo)
	test.CreateAccountBlog(t, repo, account, followedBlog)

	// Count posts (from followed blogs) that relate to python.
	count, err := find.CountSearchArticlesByAccount(account, "python")
	test.AssertNilError(t, err)

	// Should only return the three posts from followed blogs.
	test.AssertEqual(t, count, 3)
}
