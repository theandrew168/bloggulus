package ddd

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/timeutil"
)

type Session struct {
	id        uuid.UUID
	accountID uuid.UUID
	hash      string
	expiresAt time.Time

	createdAt time.Time
	updatedAt time.Time
}

type Account struct {
	id       uuid.UUID
	username string
	isAdmin  bool

	followedBlogIDs []uuid.UUID
	sessions        []Session

	createdAt time.Time
	updatedAt time.Time
}

func NewAccount(username string) (*Account, error) {
	if username == "" {
		return nil, fmt.Errorf("account: invalid username")
	}

	now := timeutil.Now()
	account := Account{
		id:       uuid.New(),
		username: username,
		isAdmin:  false,

		createdAt: now,
		updatedAt: now,
	}
	return &account, nil
}

func LoadAccount(id uuid.UUID, username string, isAdmin bool, createdAt, updatedAt time.Time) *Account {
	account := Account{
		id:       id,
		username: username,
		isAdmin:  isAdmin,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	return &account
}

func (a *Account) ID() uuid.UUID {
	return a.id
}

func (a *Account) Username() string {
	return a.username
}

func (a *Account) IsAdmin() bool {
	return a.isAdmin
}

func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Account) UpdatedAt() time.Time {
	return a.updatedAt
}

func (a *Account) FollowBlog(blogID uuid.UUID) error {
	if slices.Contains(a.followedBlogIDs, blogID) {
		return fmt.Errorf("account: already following blog")
	}

	a.followedBlogIDs = append(a.followedBlogIDs, blogID)
	return nil
}

func (a *Account) UnfollowBlog(blogID uuid.UUID) error {
	for i, id := range a.followedBlogIDs {
		if id == blogID {
			a.followedBlogIDs = slices.Delete(a.followedBlogIDs, i, i+1)
			return nil
		}
	}

	return fmt.Errorf("account: not following blog")
}

func (a *Account) Sessions() []Session {
	return a.sessions
}
