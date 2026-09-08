package webquery_test

import (
	"testing"

	webquery "github.com/theandrew168/bloggulus/backend/query/web"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestListBlogsForAccount(t *testing.T) {
	t.Parallel()

	conn, closer := test.NewDatabase(t)
	defer closer()

	repo := repository.New(conn)
	qry := webquery.NewBlog(conn)

	account := test.CreateAccount(t, repo)

	// Create and follow a blog.
	blog := test.CreateBlog(t, repo)
	account.FollowBlog(blog)

	err := repo.Account().Update(account)
	test.AssertNilError(t, err)

	// Create another blog but don't follow it.
	test.CreateBlog(t, repo)

	blogs, err := qry.ListAll(account.ID())
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
