package core

import (
	"context"
)

type AccountBlogStorage interface {
    Follow(ctx context.Context, accountID int, blogID int) error
    Unfollow(ctx context.Context, accountID int, blogID int) error
}
