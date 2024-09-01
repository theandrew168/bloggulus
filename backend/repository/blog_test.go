package repository_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestBlogCreate(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.NewBlog(t)
	err := repo.Blog().Create(blog)
	test.AssertNilError(t, err)
}

func TestBlogCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)

	// attempt to create the same blog again
	err := repo.Blog().Create(blog)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestBlogRead(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)
	got, err := repo.Blog().Read(blog.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), blog.ID())
}

func TestBlogReadByFeedURL(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)
	got, err := repo.Blog().ReadByFeedURL(blog.FeedURL())
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

	limit := 3
	offset := 0
	blogs, err := repo.Blog().List(limit, offset)
	test.AssertNilError(t, err)

	test.AssertEqual(t, len(blogs), limit)
}

func TestBlogListAll(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	test.CreateBlog(t, repo)
	test.CreateBlog(t, repo)
	test.CreateBlog(t, repo)

	blogs, err := repo.Blog().ListAll()
	test.AssertNilError(t, err)

	test.AssertAtLeast(t, len(blogs), 3)
}

func TestBlogCount(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	test.CreateBlog(t, repo)
	test.CreateBlog(t, repo)
	test.CreateBlog(t, repo)

	count, err := repo.Blog().Count()
	test.AssertNilError(t, err)

	test.AssertAtLeast(t, count, 3)
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

	err := repo.Blog().Update(blog)
	test.AssertNilError(t, err)

	got, err := repo.Blog().Read(blog.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ETag(), etag)
	test.AssertEqual(t, got.LastModified(), lastModified)
}

func TestBlogDelete(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	blog := test.CreateBlog(t, repo)

	err := repo.Blog().Delete(blog)
	test.AssertNilError(t, err)

	_, err = repo.Blog().Read(blog.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
