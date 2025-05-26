package model

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/timeutil"
)

type Account struct {
	id       uuid.UUID
	username string
	isAdmin  bool

	followedBlogIDs []uuid.UUID

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

func LoadAccount(id uuid.UUID, username string, isAdmin bool, followedBlogIDs []uuid.UUID, createdAt, updatedAt time.Time) *Account {
	account := Account{
		id:       id,
		username: username,
		isAdmin:  isAdmin,

		followedBlogIDs: followedBlogIDs,

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

func (a *Account) FollowedBlogIDs() []uuid.UUID {
	return a.followedBlogIDs
}

func (a *Account) FollowBlog(blog *Blog) error {
	if slices.Contains(a.followedBlogIDs, blog.ID()) {
		return nil
	}

	a.followedBlogIDs = append(a.followedBlogIDs, blog.ID())
	return nil
}

func (a *Account) UnfollowBlog(blog *Blog) error {
	index := slices.Index(a.followedBlogIDs, blog.ID())
	if index == -1 {
		return nil
	}

	a.followedBlogIDs = slices.Delete(a.followedBlogIDs, index, index+1)
	return nil
}

func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Account) UpdatedAt() time.Time {
	return a.updatedAt
}

func (a *Account) CheckDelete() error {
	return nil
}
