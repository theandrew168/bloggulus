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

type NewAccountParams struct {
	Username value.Name
}

func NewAccount(params NewAccountParams) (*Account, error) {
	account := Account{
		id:       uuid.New(),
		username: params.Username,
		isAdmin:  false,
		meta:     NewMeta(),
	}
	return &account, nil
}

type LoadAccountParams struct {
	ID       uuid.UUID
	Username value.Name
	IsAdmin  bool
	// TODO: Convert to a set.
	FollowedBlogIDs []uuid.UUID
	Meta            *Meta
}

func LoadAccount(params LoadAccountParams) *Account {
	account := Account{
		id:              params.ID,
		username:        params.Username,
		isAdmin:         params.IsAdmin,
		followedBlogIDs: params.FollowedBlogIDs,
		meta:            params.Meta,
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
