package query_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/query"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

func TestListArticles(t *testing.T) {
	t.Parallel()

	conn, closer := test.NewDatabase(t)
	defer closer()

	repo := repository.New(conn)
	qry := query.NewArticle(conn)

	blog := test.CreateBlog(t, repo)
	test.CreatePost(t, repo, blog)

	articles, err := qry.ListRecentArticles(1, 0)
	test.AssertNilError(t, err)

	test.AssertEqual(t, len(articles), 1)
}

func TestListArticlesByAccount(t *testing.T) {
	t.Parallel()

	conn, closer := test.NewDatabase(t)
	defer closer()

	repo := repository.New(conn)
	qry := query.NewArticle(conn)

	followedBlog := test.CreateBlog(t, repo)
	test.CreatePost(t, repo, followedBlog)
	test.CreatePost(t, repo, followedBlog)
	test.CreatePost(t, repo, followedBlog)

	unfollowedBlog := test.CreateBlog(t, repo)
	test.CreatePost(t, repo, unfollowedBlog)
	test.CreatePost(t, repo, unfollowedBlog)
	test.CreatePost(t, repo, unfollowedBlog)

	account := test.CreateAccount(t, repo)
	test.CreateAccountBlog(t, repo, account.ID(), followedBlog.ID())

	// List posts from blogs followed by this account.
	articles, err := qry.ListRecentArticlesByAccount(account.ID(), 5, 0)
	test.AssertNilError(t, err)

	// We should only get the three posts associated with the followed blog.
	test.AssertEqual(t, len(articles), 3)
}

func TestSearchArticles(t *testing.T) {
	t.Parallel()

	conn, closer := test.NewDatabase(t)
	defer closer()

	repo := repository.New(conn)
	qry := query.NewArticle(conn)

	blog := test.NewBlog(t)
	err := repo.Blog().Create(blog)
	test.AssertNilError(t, err)

	// create a post about python
	pythonPost, err := model.NewPost(model.NewPostParams{
		Blog:        blog,
		URL:         test.RandomURL(20),
		Title:       "Python",
		PublishedAt: timeutil.Now(),
		Content:     "content about python",
	})
	test.AssertNilError(t, err)

	err = repo.Post().Create(pythonPost)
	test.AssertNilError(t, err)

	// create a post about python
	boringPost, err := model.NewPost(model.NewPostParams{
		Blog:        blog,
		URL:         test.RandomURL(20),
		Title:       "Boring",
		PublishedAt: timeutil.Now(),
		Content:     "content about nothing",
	})
	test.AssertNilError(t, err)

	err = repo.Post().Create(boringPost)
	test.AssertNilError(t, err)

	// list articles that relate to python
	articles, err := qry.ListRelevantArticles("python", 1, 0)
	test.AssertNilError(t, err)

	// should find at least one
	test.AssertEqual(t, len(articles), 1)
}

func TestSearchArticlesByAccount(t *testing.T) {
	t.Parallel()

	conn, closer := test.NewDatabase(t)
	defer closer()

	repo := repository.New(conn)
	qry := query.NewArticle(conn)

	// Create some followed posts about python.
	followedBlog := test.CreateBlog(t, repo)
	for i := 0; i < 3; i++ {
		post, err := model.NewPost(model.NewPostParams{
			Blog:        followedBlog,
			URL:         test.RandomURL(20),
			Title:       "Python",
			PublishedAt: timeutil.Now(),
			Content:     "content about python",
		})
		test.AssertNilError(t, err)

		err = repo.Post().Create(post)
		test.AssertNilError(t, err)
	}

	// Create some unfollowed posts about python.
	unfollowedBlog := test.CreateBlog(t, repo)
	for i := 0; i < 3; i++ {
		post, err := model.NewPost(model.NewPostParams{
			Blog:        unfollowedBlog,
			URL:         test.RandomURL(20),
			Title:       "Python",
			PublishedAt: timeutil.Now(),
			Content:     "content about python",
		})
		test.AssertNilError(t, err)

		err = repo.Post().Create(post)
		test.AssertNilError(t, err)
	}

	account := test.CreateAccount(t, repo)
	test.CreateAccountBlog(t, repo, account.ID(), followedBlog.ID())

	// List posts (from followed blogs) that relate to python.
	articles, err := qry.ListRelevantArticlesByAccount(account.ID(), "python", 5, 0)
	test.AssertNilError(t, err)

	// Should only return the three posts from followed blogs.
	test.AssertEqual(t, len(articles), 3)
}

func TestCountArticles(t *testing.T) {
	t.Parallel()

	conn, closer := test.NewDatabase(t)
	defer closer()

	repo := repository.New(conn)
	qry := query.NewArticle(conn)

	blog := test.CreateBlog(t, repo)
	test.CreatePost(t, repo, blog)
	test.CreatePost(t, repo, blog)
	test.CreatePost(t, repo, blog)

	count, err := qry.CountRecentArticles()
	test.AssertNilError(t, err)

	test.AssertAtLeast(t, count, 3)
}

func TestCountArticlesByAccount(t *testing.T) {
	t.Parallel()

	conn, closer := test.NewDatabase(t)
	defer closer()

	repo := repository.New(conn)
	qry := query.NewArticle(conn)

	followedBlog := test.CreateBlog(t, repo)
	test.CreatePost(t, repo, followedBlog)
	test.CreatePost(t, repo, followedBlog)
	test.CreatePost(t, repo, followedBlog)

	unfollowedBlog := test.CreateBlog(t, repo)
	test.CreatePost(t, repo, unfollowedBlog)
	test.CreatePost(t, repo, unfollowedBlog)
	test.CreatePost(t, repo, unfollowedBlog)

	account := test.CreateAccount(t, repo)
	test.CreateAccountBlog(t, repo, account.ID(), followedBlog.ID())

	// We should only count the three posts associated with the followed blog.
	count, err := qry.CountRecentArticlesByAccount(account.ID())
	test.AssertNilError(t, err)
	test.AssertEqual(t, count, 3)
}

func TestCountSearchArticles(t *testing.T) {
	t.Parallel()

	conn, closer := test.NewDatabase(t)
	defer closer()

	repo := repository.New(conn)
	qry := query.NewArticle(conn)

	blog := test.CreateBlog(t, repo)

	// create a post about python
	pythonPost, err := model.NewPost(model.NewPostParams{
		Blog:        blog,
		URL:         test.RandomURL(20),
		Title:       "Python",
		PublishedAt: timeutil.Now(),
		Content:     "content about python",
	})
	test.AssertNilError(t, err)

	err = repo.Post().Create(pythonPost)
	test.AssertNilError(t, err)

	// create a post about python
	boringPost, err := model.NewPost(model.NewPostParams{
		Blog:        blog,
		URL:         test.RandomURL(20),
		Title:       "Boring",
		PublishedAt: timeutil.Now(),
		Content:     "content about nothing",
	})
	test.AssertNilError(t, err)

	err = repo.Post().Create(boringPost)
	test.AssertNilError(t, err)

	// count posts that relate to python
	count, err := qry.CountRelevantArticles("python")
	test.AssertNilError(t, err)

	// should find at least one
	test.AssertAtLeast(t, count, 1)
}

func TestCountSearchArticlesByAccount(t *testing.T) {
	t.Parallel()

	conn, closer := test.NewDatabase(t)
	defer closer()

	repo := repository.New(conn)
	qry := query.NewArticle(conn)

	// Create some followed posts about python.
	followedBlog := test.CreateBlog(t, repo)
	for i := 0; i < 3; i++ {
		post, err := model.NewPost(model.NewPostParams{
			Blog:        followedBlog,
			URL:         test.RandomURL(20),
			Title:       "Python",
			PublishedAt: timeutil.Now(),
			Content:     "content about python",
		})
		test.AssertNilError(t, err)

		err = repo.Post().Create(post)
		test.AssertNilError(t, err)
	}

	// Create some unfollowed posts about python.
	unfollowedBlog := test.CreateBlog(t, repo)
	for i := 0; i < 3; i++ {
		post, err := model.NewPost(model.NewPostParams{
			Blog:        unfollowedBlog,
			URL:         test.RandomURL(20),
			Title:       "Python",
			PublishedAt: timeutil.Now(),
			Content:     "content about python",
		})
		test.AssertNilError(t, err)

		err = repo.Post().Create(post)
		test.AssertNilError(t, err)
	}

	account := test.CreateAccount(t, repo)
	test.CreateAccountBlog(t, repo, account.ID(), followedBlog.ID())

	// Count posts (from followed blogs) that relate to python.
	count, err := qry.CountRelevantArticlesByAccount(account.ID(), "python")
	test.AssertNilError(t, err)

	// Should only return the three posts from followed blogs.
	test.AssertEqual(t, count, 3)
}
