package core_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/internal/core"
	"github.com/theandrew168/bloggulus/internal/postgresql"
)

//type Blog struct {
//	BlogID  int
//	FeedURL string
//	SiteURL string
//	Title   string
//}
//
//func NewBlog(feedURL, siteURL, title string) Blog {
//	blog := Blog{
//		FeedURL: feedURL,
//		SiteURL: siteURL,
//		Title:   title,
//	}
//	return blog
//}
//
//type BlogStorage interface {
//	Create(ctx context.Context, blog *Blog) error
//	Read(ctx context.Context, blogID int) (Blog, error)
//	ReadAll(ctx context.Context) ([]Blog, error)
//}

func TestBlogCreate(t *testing.T) {
	// check for database connection url var
	databaseURL := os.Getenv("BLOGGULUS_DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("Missing required env var: BLOGGULUS_DATABASE_URL")
	}

	// open a database connection pool
	conn, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// test connection to ensure all is well
	if err = conn.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}

	// instantiate storage interfaces
	blogStorage := postgresql.NewBlogStorage(conn)

	// generate some random data
	feedURL := "https://" + randomString(32)
	siteURL := "https://" + randomString(32)
	title := randomString(32)

	blog := core.NewBlog(feedURL, siteURL, title)
	if blog.BlogID != 0 {
		t.Fatal("blog id before creation should be zero")
	}

	// create an example blog
	err = blogStorage.Create(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	if blog.BlogID == 0 {
		t.Fatal("blog id after creation should be nonzero")
	}
}
