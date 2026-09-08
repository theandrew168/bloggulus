package test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/timeutil"
	"github.com/theandrew168/bloggulus/backend/value"
)

func NewCount(count int) value.Count {
	c, err := value.NewCount(count)
	if err != nil {
		panic(err)
	}

	return c
}

func NewName(name string) value.Name {
	n, err := value.NewName(name)
	if err != nil {
		panic(err)
	}

	return n
}

func NewURL(url string) value.URL {
	u, err := value.NewURL(url)
	if err != nil {
		panic(err)
	}

	return u
}

func RandomString(n int) string {
	valid := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	buf := make([]byte, n)
	for i := range buf {
		buf[i] = valid[rand.Intn(len(valid))]
	}

	return string(buf)
}

func RandomTime() time.Time {
	return timeutil.Now()
}

func RandomName(n int) value.Name {
	return NewName(RandomString(n))
}

func RandomURL(n int) value.URL {
	return NewURL("https://" + RandomString(n))
}

func NewBlogParams() model.NewBlogParams {
	return model.NewBlogParams{
		FeedURL:      RandomURL(32),
		SiteURL:      RandomURL(32),
		Title:        RandomName(32),
		SyncedAt:     RandomTime(),
		ETag:         RandomString(32),
		LastModified: RandomString(32),
	}
}

func NewPostParams(blog *model.Blog) model.NewPostParams {
	return model.NewPostParams{
		Blog:        blog,
		URL:         RandomURL(32),
		Title:       RandomName(32),
		PublishedAt: RandomTime(),
		Content:     RandomString(32),
	}
}

func NewTagParams() model.NewTagParams {
	return model.NewTagParams{
		Name: RandomName(32),
	}
}

func NewAccountParams() model.NewAccountParams {
	return model.NewAccountParams{
		Username: RandomName(32),
	}
}

func NewSessionParams(account *model.Account) model.NewSessionParams {
	return model.NewSessionParams{
		Account: account,
		TTL:     24 * time.Hour,
	}
}

func NewBlog() *model.Blog {
	blog, err := model.NewBlog(NewBlogParams())
	if err != nil {
		panic(err)
	}

	// TODO: Update tests to account for visibility and then remove this.
	blog.SetIsPublic(true)
	return blog
}

func NewPost(blog *model.Blog) *model.Post {
	post, err := model.NewPost(NewPostParams(blog))
	if err != nil {
		panic(err)
	}

	return post
}

func NewTag() *model.Tag {
	tag, err := model.NewTag(NewTagParams())
	if err != nil {
		panic(err)
	}

	return tag
}

func NewAccount() *model.Account {
	account, err := model.NewAccount(NewAccountParams())
	if err != nil {
		panic(err)
	}

	return account
}

func NewSession(account *model.Account) (*model.Session, value.Token) {
	session, sessionToken, err := model.NewSession(NewSessionParams(account))
	if err != nil {
		panic(err)
	}

	return session, sessionToken
}

// mocks a blog and creates it in the database
func CreateBlog(t *testing.T, repo *repository.Repository) *model.Blog {
	t.Helper()

	// generate some random blog data
	blog := NewBlog()

	// create an example blog
	err := repo.Blog().Create(blog)
	AssertNilError(t, err)

	return blog
}

// mocks a post and creates it in the database
func CreatePost(t *testing.T, repo *repository.Repository, blog *model.Blog) *model.Post {
	t.Helper()

	// generate some random post data
	post := NewPost(blog)

	// create an example post
	err := repo.Post().Create(post)
	AssertNilError(t, err)

	return post
}

// mocks a tag and creates it in the database
func CreateTag(t *testing.T, repo *repository.Repository) *model.Tag {
	t.Helper()

	// generate some random tag data
	tag := NewTag()

	// create an example tag
	err := repo.Tag().Create(tag)
	AssertNilError(t, err)

	return tag
}

// mocks an account and creates it in the database
func CreateAccount(t *testing.T, repo *repository.Repository) *model.Account {
	t.Helper()

	// generate some random account data
	account := NewAccount()

	// create an example account
	err := repo.Account().Create(account)
	AssertNilError(t, err)

	return account
}

// mocks a session and creates it in the database
func CreateSession(t *testing.T, repo *repository.Repository, account *model.Account) (*model.Session, value.Token) {
	t.Helper()

	// generate some random session data
	session, sessionToken := NewSession(account)

	// create an example session
	err := repo.Session().Create(session)
	AssertNilError(t, err)

	return session, sessionToken
}
