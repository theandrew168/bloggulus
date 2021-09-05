package core

import (
	"context"
)

type FollowStorage interface {
	Follow(ctx context.Context, accountID int, blogID int) error
	Unfollow(ctx context.Context, accountID int, blogID int) error
}
