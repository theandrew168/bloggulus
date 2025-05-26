package command

import (
	"errors"

	"github.com/theandrew168/bloggulus/backend/command/sync"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

// Sync a new or existing Blog based on the provided feed URL.
func (cmd *Command) SyncBlog(feedURL string) error {
	blog, err := cmd.repo.Blog().ReadByFeedURL(feedURL)
	if err != nil {
		if !errors.Is(err, postgres.ErrNotFound) {
			return err
		}

		// An ErrNotFound is acceptable (and expected) here. The only difference
		// is that we won't be able to include the ETag and Last-Modified headers
		// in the request. This is fine for new blogs (an unconditional fetch).
		return sync.SyncNewBlog(cmd.repo, cmd.feedFetcher, feedURL)
	}

	return sync.SyncExistingBlog(cmd.repo, cmd.feedFetcher, blog)
}
