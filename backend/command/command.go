package command

import (
	"github.com/theandrew168/bloggulus/backend/repository"
)

type Command struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Command {
	cmd := Command{
		repo: repo,
	}
	return &cmd
}
