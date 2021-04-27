package storage

import (
	"context"
)

type AccountBlog interface {
	Follow(ctx context.Context, accountID int, blogID int) error
	Unfollow(ctx context.Context, accountID int, blogID int) error
}
