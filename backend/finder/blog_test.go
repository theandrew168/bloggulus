package finder_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/test"
)

func TestListBlogsForAccount(t *testing.T) {
	t.Parallel()

	repo, repoCloser := test.NewRepository(t)
	defer repoCloser()

	find, findCloser := test.NewFinder(t)
	defer findCloser()

	account := test.CreateAccount(t, repo)

	// Create and follow a blog.
	blog := test.CreateBlog(t, repo)
	test.CreateAccountBlog(t, repo, account, blog)

	// Create another blog but don't follow it.
	test.CreateBlog(t, repo)

	blogs, err := find.ListBlogsForAccount(account)
	test.AssertNilError(t, err)

	// Count how many blogs are being followed.
	followed := 0
	for _, b := range blogs {
		if b.IsFollowing {
			followed += 1
		}
	}

	// Should only be one.
	test.AssertEqual(t, followed, 1)
}
