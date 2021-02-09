package models

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func TestSourcedPostStorage(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("Missing required env var: DATABASE_URL")
	}

	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err = Migrate(db, "migrations/*.sql"); err != nil {
		t.Fatal(err)
	}

	blogStorage := NewBlogStorage(db)
	blog, err := blogStorage.Create(ctx, "https://example.com/atom.xml", "https://example.com", "FooBar")
	if err != nil {
		t.Fatal(err)
	}

	postStorage := NewPostStorage(db)
	post, err := postStorage.Create(ctx, blog.BlogID, "https://example.com/blog/1", "Blog 1", time.Now())
	if err != nil {
		t.Fatal(err)
	}

	sourcedPostStorage := NewSourcedPostStorage(db)
	posts, err := sourcedPostStorage.ReadRecent(ctx, 20)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) == 0 {
		t.Fatal("expected at least one post")
	}

	if err = postStorage.Delete(ctx, post.PostID); err != nil {
		t.Fatal(err)
	}

	if err = blogStorage.Delete(ctx, blog.BlogID); err != nil {
		t.Fatal(err)
	}
}
