package repository_test

import (
	"context"
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestBlogCreate(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.NewBlog()
	err := repo.Blog().Create(context.Background(), blog)
	test.AssertNilError(t, err)
}

func TestBlogCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)

	// attempt to create the same blog again
	err := repo.Blog().Create(context.Background(), blog)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestBlogRead(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)
	got, err := repo.Blog().Read(context.Background(), blog.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), blog.ID())
}

func TestBlogReadByFeedURL(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)
	got, err := repo.Blog().ReadByFeedURL(context.Background(), blog.FeedURL())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), blog.ID())
}

func TestBlogList(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	test.CreateBlog(t, repo)
	test.CreateBlog(t, repo)
	test.CreateBlog(t, repo)

	blogs, err := repo.Blog().List(context.Background())
	test.AssertNilError(t, err)

	test.AssertAtLeast(t, len(blogs), 3)
}

func TestBlogUpdate(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)

	etag := "foo"
	blog.SetETag(etag)

	lastModified := "bar"
	blog.SetLastModified(lastModified)

	err := repo.Blog().Update(context.Background(), blog)
	test.AssertNilError(t, err)

	got, err := repo.Blog().Read(context.Background(), blog.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ETag(), etag)
	test.AssertEqual(t, got.LastModified(), lastModified)
}

func TestBlogDelete(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)

	err := repo.Blog().Delete(context.Background(), blog)
	test.AssertNilError(t, err)

	_, err = repo.Blog().Read(context.Background(), blog.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
