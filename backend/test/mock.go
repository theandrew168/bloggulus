package test

import (
	"testing"
	"time"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/storage"
)

func NewBlog() *admin.Blog {
	blog := admin.NewBlog(
		RandomURL(32),
		RandomURL(32),
		RandomString(32),
		RandomString(32),
		RandomString(32),
		RandomTime(),
	)
	return blog
}

func NewPost(blog *admin.Blog) *admin.Post {
	post := admin.NewPost(
		blog,
		RandomURL(32),
		RandomString(32),
		RandomString(32),
		RandomTime(),
	)
	return post
}

func NewTag() *admin.Tag {
	tag := admin.NewTag(
		RandomString(32),
	)
	return tag
}

func NewAccount(t *testing.T) *admin.Account {
	account, err := admin.NewAccount(
		RandomString(32),
		RandomString(32),
	)
	AssertNilError(t, err)

	return account
}

func NewToken(t *testing.T, account *admin.Account) *admin.Token {
	token, err := admin.NewToken(
		account,
		RandomString(32),
		// expire in 24 hours
		time.Now().UTC().Add(24*time.Hour),
	)
	AssertNilError(t, err)

	return token
}

// mocks a blog and creates it in the database
func CreateBlog(t *testing.T, store *storage.Storage) *admin.Blog {
	t.Helper()

	// generate some random blog data
	blog := NewBlog()

	// create an example blog
	err := store.Admin().Blog().Create(blog)
	AssertNilError(t, err)

	return blog
}

// mocks a post and creates it in the database
func CreatePost(t *testing.T, store *storage.Storage) *admin.Post {
	t.Helper()

	// generate some random blog data
	blog := NewBlog()

	// create an example blog
	err := store.Admin().Blog().Create(blog)
	AssertNilError(t, err)

	// generate some random post data
	post := NewPost(blog)

	// create an example post
	err = store.Admin().Post().Create(post)
	AssertNilError(t, err)

	return post
}

// mocks a tag and creates it in the database
func CreateTag(t *testing.T, store *storage.Storage) *admin.Tag {
	t.Helper()

	// generate some random tag data
	tag := NewTag()

	// create an example tag
	err := store.Admin().Tag().Create(tag)
	AssertNilError(t, err)

	return tag
}

// mocks an account and creates it in the database
func CreateAccount(t *testing.T, store *storage.Storage) *admin.Account {
	t.Helper()

	// generate some random account data
	account := NewAccount(t)

	// create an example account
	err := store.Admin().Account().Create(account)
	AssertNilError(t, err)

	return account
}

// mocks a token and creates it in the database
func CreateToken(t *testing.T, store *storage.Storage) *admin.Token {
	t.Helper()

	// generate some random account data
	account := NewAccount(t)

	// create an example account
	err := store.Admin().Account().Create(account)
	AssertNilError(t, err)

	// generate some random token data
	token := NewToken(t, account)

	// create an example token
	err = store.Admin().Token().Create(token)
	AssertNilError(t, err)

	return token
}
