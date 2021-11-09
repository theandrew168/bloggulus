package feed

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"

	"github.com/theandrew168/bloggulus/internal/core"
)

type Reader interface {
	ReadBlog(feedURL string) (core.Blog, error)
	ReadBlogPosts(blog core.Blog) ([]core.Post, error)
	ReadPostBody(post core.Post) (string, error)
}

// TODO: pull CleanHTML into pure func

// I know...
var (
	codePattern   = regexp.MustCompile(`(?s)<code>.*?</code>`)
	footerPattern = regexp.MustCompile(`(?s)<footer>.*?</footer>`)
	headerPattern = regexp.MustCompile(`(?s)<header>.*?</header>`)
	navPattern    = regexp.MustCompile(`(?s)<nav>.*?</nav>`)
	prePattern    = regexp.MustCompile(`(?s)<pre>.*?</pre>`)
)

type reader struct{}

func NewReader() Reader {
	r := reader{}
	return &r
}

func (r *reader) ReadBlog(feedURL string) (core.Blog, error) {
	// early check to ensure the URL is valid
	_, err := url.Parse(feedURL)
	if err != nil {
		return core.Blog{}, err
	}

	// attempt to parse the feed via gofeed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return core.Blog{}, err
	}

	// create a core.Blog for the feed
	blog := core.NewBlog(feedURL, feed.Link, feed.Title)
	return blog, nil
}

func (r *reader) ReadBlogPosts(blog core.Blog) ([]core.Post, error) {
	// attempt to parse the feed via gofeed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(blog.FeedURL)
	if err != nil {
		return nil, err
	}

	// create a core.Post for each entry
	var posts []core.Post
	for _, item := range feed.Items {
		// try Updated then Published to obtain a timestamp
		var updated time.Time
		if item.UpdatedParsed != nil {
			updated = *item.UpdatedParsed
		} else if item.PublishedParsed != nil {
			updated = *item.PublishedParsed
		} else {
			// else default to now
			updated = time.Now()
		}

		post := core.NewPost(item.Link, item.Title, updated, blog)
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *reader) ReadPostBody(post core.Post) (string, error) {
	resp, err := http.Get(post.URL)
	if err != nil {
		return "", fmt.Errorf("%v: %v", post.URL, err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%v: %v", post.URL, err)
	}

	// please PR a better way :(
	buf = codePattern.ReplaceAll(buf, nil)
	buf = footerPattern.ReplaceAll(buf, nil)
	buf = headerPattern.ReplaceAll(buf, nil)
	buf = navPattern.ReplaceAll(buf, nil)
	buf = prePattern.ReplaceAll(buf, nil)

	p := bluemonday.StrictPolicy()
	body := string(buf)
	body = p.Sanitize(body)
	body = html.UnescapeString(body)
	body = strings.ToValidUTF8(body, "")

	return body, nil
}

type mockReader struct {
	blog  core.Blog
	posts []core.Post
	body  string
}

func NewMockReader(blog core.Blog, posts []core.Post, body string) Reader {
	r := mockReader{
		blog:  blog,
		posts: posts,
		body:  body,
	}
	return &r
}

func (r *mockReader) ReadBlog(feedURL string) (core.Blog, error) {
	return r.blog, nil
}

func (r *mockReader) ReadBlogPosts(blog core.Blog) ([]core.Post, error) {
	return r.posts, nil
}

func (r *mockReader) ReadPostBody(post core.Post) (string, error) {
	return r.body, nil
}
