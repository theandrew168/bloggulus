package mock

import (
	"encoding/xml"
	"time"

	"github.com/theandrew168/bloggulus/backend/domain/admin/feed"
)

type xmlLink struct {
	HREF string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
}

type xmlPost struct {
	URL         xmlLink   `xml:"link"`
	Title       string    `xml:"title"`
	Content     string    `xml:"content"`
	PublishedAt time.Time `xml:"published"`
}

type xmlBlog struct {
	XMLName xml.Name  `xml:"feed"`
	Links   []xmlLink `xml:"link"`
	Title   string    `xml:"title"`
	Posts   []xmlPost `xml:"entry"`
}

// Convert a feed.Blog into an XML (Atom) document.
func GenerateAtomFeed(blog feed.Blog) (string, error) {
	var posts []xmlPost
	for _, post := range blog.Posts {
		posts = append(posts, xmlPost{
			URL:         xmlLink{HREF: post.URL},
			Title:       post.Title,
			Content:     post.Content,
			PublishedAt: post.PublishedAt,
		})
	}

	b := xmlBlog{
		Links: []xmlLink{
			{HREF: blog.FeedURL, Rel: "self"},
			{HREF: blog.SiteURL, Rel: "alternate"},
		},
		Title: blog.Title,
		Posts: posts,
	}

	out, err := xml.Marshal(b)
	if err != nil {
		return "", err
	}

	return xml.Header + string(out), nil
}
