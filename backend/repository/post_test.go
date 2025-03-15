package repository_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestPostCreate(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.NewBlog(t)
	err := repo.Blog().Create(blog)
	test.AssertNilError(t, err)

	post := test.NewPost(t, blog)
	err = repo.Post().Create(post)
	test.AssertNilError(t, err)
}

func TestPostCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)
	post := test.CreatePost(t, repo, blog)

	// attempt to create the same post again
	err := repo.Post().Create(post)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestPostRead(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)
	post := test.CreatePost(t, repo, blog)
	got, err := repo.Post().Read(post.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), post.ID())
}

func TestPostReadByURL(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)
	post := test.CreatePost(t, repo, blog)
	got, err := repo.Post().ReadByURL(post.URL())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), post.ID())
}

func TestPostListByBlog(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)
	test.CreatePost(t, repo, blog)
	test.CreatePost(t, repo, blog)
	test.CreatePost(t, repo, blog)

	posts, err := repo.Post().ListByBlog(blog)
	test.AssertNilError(t, err)

	test.AssertEqual(t, len(posts), 3)
}

func TestPostCountByBlog(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)
	test.CreatePost(t, repo, blog)
	test.CreatePost(t, repo, blog)
	test.CreatePost(t, repo, blog)

	count, err := repo.Post().CountByBlog(blog)
	test.AssertNilError(t, err)

	test.AssertEqual(t, count, 3)
}

func TestPostUpdate(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)
	post := test.CreatePost(t, repo, blog)

	content := "foobar"
	post.SetContent(content)

	err := repo.Post().Update(post)
	test.AssertNilError(t, err)

	got, err := repo.Post().Read(post.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.Content(), content)
}

func TestPostDelete(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)
	post := test.CreatePost(t, repo, blog)

	err := repo.Post().Delete(post)
	test.AssertNilError(t, err)

	_, err = repo.Post().Read(post.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
