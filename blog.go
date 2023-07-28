package bloggulus

type Blog struct {
	FeedURL      string `json:"feed_url"`
	SiteURL      string `json:"site_url"`
	Title        string `json:"title"`
	ETag         string `json:"etag"`
	LastModified string `json:"last_modified"`

	// readonly (from database, after creation)
	ID int `json:"id"`
}

// NewBlog creates a new Blog struct.
func NewBlog(feedURL, siteURL, title, etag, lastModified string) Blog {
	blog := Blog{
		FeedURL:      feedURL,
		SiteURL:      siteURL,
		Title:        title,
		ETag:         etag,
		LastModified: lastModified,
	}
	return blog
}
