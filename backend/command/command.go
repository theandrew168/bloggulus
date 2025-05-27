package command

import (
	"github.com/theandrew168/bloggulus/backend/feed"
	"github.com/theandrew168/bloggulus/backend/repository"
)

type Command struct {
	repo        *repository.Repository
	feedFetcher feed.FeedFetcher
}

func New(repo *repository.Repository, feedFetcher feed.FeedFetcher) *Command {
	cmd := Command{
		repo:        repo,
		feedFetcher: feedFetcher,
	}
	return &cmd
}
