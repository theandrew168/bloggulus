package test

import (
	"testing"
	"time"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/value"
)

func NewBlog(t *testing.T) *model.Blog {
	blog, err := model.NewBlog(model.NewBlogParams{
		FeedURL:      RandomURL(32),
		SiteURL:      RandomURL(32),
		Title:        RandomString(32),
		SyncedAt:     RandomTime(),
		ETag:         RandomString(32),
		LastModified: RandomString(32),
	})
	AssertNilError(t, err)

	return blog
}

func NewPost(t *testing.T, blog *model.Blog) *model.Post {
	post, err := model.NewPost(model.NewPostParams{
		Blog:        blog,
		URL:         RandomURL(32),
		Title:       RandomString(32),
		PublishedAt: RandomTime(),
		Content:     RandomString(32),
	})
	AssertNilError(t, err)

	return post
}

func NewTag(t *testing.T) *model.Tag {
	tag, err := model.NewTag(model.NewTagParams{
		Name: RandomName(t),
	})
	AssertNilError(t, err)

	return tag
}

func NewAccount(t *testing.T) *model.Account {
	account, err := model.NewAccount(model.NewAccountParams{
		Username: RandomName(t),
	})
	AssertNilError(t, err)

	return account
}

func NewSession(t *testing.T, account *model.Account) (*model.Session, value.Token) {
	session, sessionToken, err := model.NewSession(model.NewSessionParams{
		Account: account,
		TTL:     24 * time.Hour,
	})
	AssertNilError(t, err)

	return session, sessionToken
}

// mocks a blog and creates it in the database
func CreateBlog(t *testing.T, repo *repository.Repository) *model.Blog {
	t.Helper()

	// generate some random blog data
	blog := NewBlog(t)

	// create an example blog
	err := repo.Blog().Create(blog)
	AssertNilError(t, err)

	return blog
}

// mocks a post and creates it in the database
func CreatePost(t *testing.T, repo *repository.Repository, blog *model.Blog) *model.Post {
	t.Helper()

	// generate some random post data
	post := NewPost(t, blog)

	// create an example post
	err := repo.Post().Create(post)
	AssertNilError(t, err)

	return post
}

// mocks a tag and creates it in the database
func CreateTag(t *testing.T, repo *repository.Repository) *model.Tag {
	t.Helper()

	// generate some random tag data
	tag := NewTag(t)

	// create an example tag
	err := repo.Tag().Create(tag)
	AssertNilError(t, err)

	return tag
}

// mocks an account and creates it in the database
func CreateAccount(t *testing.T, repo *repository.Repository) *model.Account {
	t.Helper()

	// generate some random account data
	account := NewAccount(t)

	// create an example account
	err := repo.Account().Create(account)
	AssertNilError(t, err)

	return account
}

// mocks a session and creates it in the database
func CreateSession(t *testing.T, repo *repository.Repository, account *model.Account) (*model.Session, value.Token) {
	t.Helper()

	// generate some random session data
	session, sessionToken := NewSession(t, account)

	// create an example session
	err := repo.Session().Create(session)
	AssertNilError(t, err)

	return session, sessionToken
}

// create an account blog in the database
func CreateAccountBlog(t *testing.T, repo *repository.Repository, account *model.Account, blog *model.Blog) {
	t.Helper()

	// create an account blog
	err := repo.AccountBlog().Create(account, blog)
	AssertNilError(t, err)
}
