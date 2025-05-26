package command

import (
	"github.com/theandrew168/bloggulus/backend/fetch"
	"github.com/theandrew168/bloggulus/backend/repository"
)

type Command struct {
	repo        *repository.Repository
	feedFetcher fetch.FeedFetcher
}

func New(repo *repository.Repository, feedFetcher fetch.FeedFetcher) *Command {
	cmd := Command{
		repo:        repo,
		feedFetcher: feedFetcher,
	}
	return &cmd
}
