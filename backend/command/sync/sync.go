package sync

import (
	"log/slog"

	"github.com/theandrew168/bloggulus/backend/feed"
	"github.com/theandrew168/bloggulus/backend/fetch"
	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

// UpdateCacheHeaders updates the ETag and Last-Modified headers for a blog if they have changed.
func UpdateCacheHeaders(blog *model.Blog, response fetch.FetchFeedResponse) bool {
	headersChanged := false
	if response.ETag != "" && response.ETag != blog.ETag() {
		headersChanged = true
		blog.SetETag(response.ETag)
	}

	if response.LastModified != "" && response.LastModified != blog.LastModified() {
		headersChanged = true
		blog.SetLastModified(response.LastModified)
	}

	return headersChanged
}

type ComparePostsResult struct {
	PostsToCreate []*model.Post
	PostsToUpdate []*model.Post
}

// ComparePosts compares a list of known posts to a list of feed posts and returns
// a list of posts to create and a list of posts to update.
func ComparePosts(blog *model.Blog, knownPosts []*model.Post, feedPosts []feed.Post) (ComparePostsResult, error) {
	// Create a map of URLs to posts for quick lookups.
	knownPostsByURL := make(map[string]*model.Post)
	for _, post := range knownPosts {
		knownPostsByURL[post.URL()] = post
	}

	var postsToCreate []*model.Post
	var postsToUpdate []*model.Post

	// Compare each post in the feed to the posts in the database.
	for _, feedPost := range feedPosts {
		knownPost, ok := knownPostsByURL[feedPost.URL]
		if !ok {
			// The post is new so we need to create it.
			postToCreate, err := model.NewPost(
				blog,
				feedPost.URL,
				feedPost.Title,
				feedPost.Content,
				feedPost.PublishedAt,
			)
			if err != nil {
				return ComparePostsResult{}, err
			}

			postsToCreate = append(postsToCreate, postToCreate)
		} else {
			// The post already exists but we might need to update it.
			knownPostShouldBeUpdated := false

			// Check if the post's title has changed.
			if feedPost.Title != "" && feedPost.Title != knownPost.Title() {
				knownPost.SetTitle(feedPost.Title)
				knownPostShouldBeUpdated = true
			}

			// Check if the post's content has changed.
			if feedPost.Content != "" && feedPost.Content != knownPost.Content() {
				knownPost.SetContent(feedPost.Content)
				knownPostShouldBeUpdated = true
			}

			// Check if the post's publishedAt date has changed.
			if feedPost.PublishedAt != knownPost.PublishedAt() {
				knownPost.SetPublishedAt(feedPost.PublishedAt)
				knownPostShouldBeUpdated = true
			}

			// If any post data has changed, add it to the list of posts to update.
			if knownPostShouldBeUpdated {
				postsToUpdate = append(postsToUpdate, knownPost)
			}
		}
	}

	result := ComparePostsResult{
		PostsToCreate: postsToCreate,
		PostsToUpdate: postsToUpdate,
	}
	return result, nil
}

func SyncNewBlog(repo *repository.Repository, feedFetcher fetch.FeedFetcher, feedURL string) error {
	// Make an unconditional fetch for the blog's feed.
	req := fetch.FetchFeedRequest{
		URL: feedURL,
	}
	resp, err := feedFetcher.FetchFeed(req)
	if err != nil {
		return err
	}

	// No feed data from a new blog is an error.
	if resp.Feed == "" {
		return fetch.ErrUnreachableFeed
	}

	feedBlog, err := feed.Parse(feedURL, resp.Feed)
	if err != nil {
		return err
	}

	// Create a new blog based on the feed data.
	blog, err := model.NewBlog(
		feedBlog.FeedURL,
		feedBlog.SiteURL,
		feedBlog.Title,
		resp.ETag,
		resp.LastModified,
		timeutil.Now(),
	)
	if err != nil {
		return err
	}

	err = repo.Blog().Create(blog)
	if err != nil {
		return err
	}

	err = SyncPosts(repo, blog, feedBlog.Posts)
	if err != nil {
		return err
	}

	return nil
}

func SyncExistingBlog(repo *repository.Repository, feedFetcher fetch.FeedFetcher, blog *model.Blog) error {
	// Make a conditional fetch for the blog's feed.
	req := fetch.FetchFeedRequest{
		URL:          blog.FeedURL(),
		ETag:         blog.ETag(),
		LastModified: blog.LastModified(),
	}
	resp, err := feedFetcher.FetchFeed(req)
	if err != nil {
		return err
	}

	// Update the blog's cache headers if they have changed.
	headersChanged := UpdateCacheHeaders(blog, resp)
	if headersChanged {
		err = repo.Blog().Update(blog)
		if err != nil {
			return err
		}
	}

	if resp.Feed == "" {
		slog.Info("skipping blog (no feed content)", "title", blog.Title(), "id", blog.ID())
		return nil
	}

	feedBlog, err := feed.Parse(blog.FeedURL(), resp.Feed)
	if err != nil {
		return err
	}

	err = SyncPosts(repo, blog, feedBlog.Posts)
	if err != nil {
		return err
	}

	return nil
}

func SyncPosts(repo *repository.Repository, blog *model.Blog, feedPosts []feed.Post) error {
	// List all known posts for the current blog.
	knownPosts, err := repo.Post().ListByBlog(blog)
	if err != nil {
		return err
	}

	// Compare the known posts to the feed posts.
	result, err := ComparePosts(blog, knownPosts, feedPosts)
	if err != nil {
		return err
	}

	// Create any posts that are new.
	for _, post := range result.PostsToCreate {
		err = repo.Post().Create(post)
		if err != nil {
			slog.Warn("failed to create post", "url", post.URL(), "error", err.Error())
		}
	}

	// Update any posts that have changed.
	for _, post := range result.PostsToUpdate {
		err = repo.Post().Update(post)
		if err != nil {
			slog.Warn("failed to update post", "url", post.URL(), "error", err.Error())
		}
	}

	return nil
}
