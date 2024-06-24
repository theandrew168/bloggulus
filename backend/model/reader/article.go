package reader

import (
	"time"
)

type Article struct {
	title       string
	url         string
	blogTitle   string
	blogURL     string
	publishedAt time.Time
	tags        []string
}

func LoadArticle(title, url, blogTitle, blogURL string, publishedAt time.Time, tags []string) *Article {
	post := Article{
		title:       title,
		url:         url,
		blogTitle:   blogTitle,
		blogURL:     blogURL,
		publishedAt: publishedAt,
		tags:        tags,
	}
	return &post
}

func (a *Article) Title() string {
	return a.title
}

func (a *Article) URL() string {
	return a.url
}

func (a *Article) BlogTitle() string {
	return a.blogTitle
}

func (a *Article) BlogURL() string {
	return a.blogURL
}

func (a *Article) PublishedAt() time.Time {
	return a.publishedAt
}

func (a *Article) Tags() []string {
	return a.tags
}
