package models

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

func TestBlogStorage(t *testing.T) {
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

	_, err = blogStorage.Read(ctx, blog.BlogID)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: check values?

	blogs, err := blogStorage.ReadAll(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(blogs) == 0 {
		t.Fatal("expected at least one blog")
	}

	if err = blogStorage.Delete(ctx, blog.BlogID); err != nil {
		t.Fatal(err)
	}
}
