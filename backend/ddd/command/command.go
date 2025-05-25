package command

import "github.com/google/uuid"

type Command interface {
	Kind() string
}

type AddBlog struct {
	URL string
}

func (c AddBlog) Kind() string {
	return "AddBlog"
}

type FollowBlog struct {
	AccountID uuid.UUID
	BlogID    uuid.UUID
}

func (c FollowBlog) Kind() string {
	return "FollowBlog"
}

type UnfollowBlog struct {
	AccountID uuid.UUID
	BlogID    uuid.UUID
}

func (c UnfollowBlog) Kind() string {
	return "UnfollowBlog"
}
