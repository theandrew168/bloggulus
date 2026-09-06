package command

import (
	"github.com/theandrew168/bloggulus/backend/feed"
	"github.com/theandrew168/bloggulus/backend/repository"
)

type Command struct {
	account *AccountCommand
	auth    *AuthCommand
	blog    *BlogCommand
	sync    *SyncCommand
}

func New(repo *repository.Repository, feedFetcher feed.FeedFetcher) *Command {
	cmd := Command{
		account: NewAccount(repo),
		auth:    NewAuth(repo),
		blog:    NewBlog(repo),
		sync:    NewSync(repo, feedFetcher),
	}
	return &cmd
}

func (cmd *Command) Account() *AccountCommand {
	return cmd.account
}

func (cmd *Command) Auth() *AuthCommand {
	return cmd.auth
}

func (cmd *Command) Blog() *BlogCommand {
	return cmd.blog
}

func (cmd *Command) Sync() *SyncCommand {
	return cmd.sync
}
