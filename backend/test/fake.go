package test

import (
	"testing"
	"time"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/repository"
)

func NewBlog(t *testing.T) *model.Blog {
	blog, err := model.NewBlog(
		RandomURL(32),
		RandomURL(32),
		RandomString(32),
		RandomString(32),
		RandomString(32),
		RandomTime(),
	)
	AssertNilError(t, err)

	return blog
}

func NewPost(t *testing.T, blog *model.Blog) *model.Post {
	post, err := model.NewPost(
		blog,
		RandomURL(32),
		RandomString(32),
		RandomString(32),
		RandomTime(),
	)
	AssertNilError(t, err)

	return post
}

func NewTag(t *testing.T) *model.Tag {
	tag, err := model.NewTag(
		RandomString(32),
	)
	AssertNilError(t, err)

	return tag
}

func NewAccount(t *testing.T) (*model.Account, string) {
	password := RandomString(32)
	account, err := model.NewAccount(
		RandomString(32),
		password,
	)
	AssertNilError(t, err)

	return account, password
}

func NewSession(t *testing.T, account *model.Account) (*model.Session, string) {
	session, sessionID, err := model.NewSession(
		account,
		// expire in 24 hours
		24*time.Hour,
	)
	AssertNilError(t, err)

	return session, sessionID
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
func CreateAccount(t *testing.T, repo *repository.Repository) (*model.Account, string) {
	t.Helper()

	// generate some random account data
	account, password := NewAccount(t)

	// create an example account
	err := repo.Account().Create(account)
	AssertNilError(t, err)

	return account, password
}

// mocks a session and creates it in the database
func CreateSession(t *testing.T, repo *repository.Repository, account *model.Account) (*model.Session, string) {
	t.Helper()

	// generate some random session data
	session, sessionID := NewSession(t, account)

	// create an example session
	err := repo.Session().Create(session)
	AssertNilError(t, err)

	return session, sessionID
}

// create an account blog in the database
func CreateAccountBlog(t *testing.T, repo *repository.Repository, account *model.Account, blog *model.Blog) {
	t.Helper()

	// create an account blog
	err := repo.AccountBlog().Create(account, blog)
	AssertNilError(t, err)
}
