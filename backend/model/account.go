package model

import (
	"slices"
	"uuid"

	"github.com/theandrew168/bloggulus/backend/value"
)

type Account struct {
	id       uuid.UUID
	username value.Name
	isAdmin  bool

	followedBlogIDs []uuid.UUID

	meta *Meta
}

func NewAccount(username value.Name) (*Account, error) {
	account := Account{
		id:       uuid.New(),
		username: username,
		isAdmin:  false,
		meta:     NewMeta(),
	}
	return &account, nil
}

func LoadAccount(id uuid.UUID, username value.Name, isAdmin bool, followedBlogIDs []uuid.UUID, meta *Meta) *Account {
	account := Account{
		id:       id,
		username: username,
		isAdmin:  isAdmin,

		followedBlogIDs: followedBlogIDs,

		meta: meta,
	}
	return &account
}

func (a *Account) ID() uuid.UUID {
	return a.id
}

func (a *Account) Username() value.Name {
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

func (a *Account) Meta() *Meta {
	return a.meta
}
