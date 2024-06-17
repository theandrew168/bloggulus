package reader

import (
	"time"
)

type Post struct {
	title       string
	url         string
	blogTitle   string
	blogURL     string
	publishedAt time.Time
	tags        []string
}

func LoadPost(title, url, blogTitle, blogURL string, publishedAt time.Time, tags []string) *Post {
	post := Post{
		title:       title,
		url:         url,
		blogTitle:   blogTitle,
		blogURL:     blogURL,
		publishedAt: publishedAt,
		tags:        tags,
	}
	return &post
}

func (p *Post) Title() string {
	return p.title
}

func (p *Post) URL() string {
	return p.url
}

func (p *Post) BlogTitle() string {
	return p.blogTitle
}

func (p *Post) BlogURL() string {
	return p.blogURL
}

func (p *Post) PublishedAt() time.Time {
	return p.publishedAt
}

func (p *Post) Tags() []string {
	return p.tags
}
