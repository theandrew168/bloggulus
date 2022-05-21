package feed

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"

	"github.com/theandrew168/bloggulus"
)

// I know...
var (
	codePattern   = regexp.MustCompile(`(?s)<code>.*?</code>`)
	footerPattern = regexp.MustCompile(`(?s)<footer>.*?</footer>`)
	headerPattern = regexp.MustCompile(`(?s)<header>.*?</header>`)
	navPattern    = regexp.MustCompile(`(?s)<nav>.*?</nav>`)
	prePattern    = regexp.MustCompile(`(?s)<pre>.*?</pre>`)
)

// please PR a better way :(
func CleanHTML(s string) string {
	s = codePattern.ReplaceAllString(s, "")
	s = footerPattern.ReplaceAllString(s, "")
	s = headerPattern.ReplaceAllString(s, "")
	s = navPattern.ReplaceAllString(s, "")
	s = prePattern.ReplaceAllString(s, "")

	s = bluemonday.StrictPolicy().Sanitize(s)
	s = html.UnescapeString(s)
	s = strings.ToValidUTF8(s, "")

	return s
}

type Reader interface {
	ReadBlog(feedURL string) (bloggulus.Blog, error)
	ReadBlogPosts(blog bloggulus.Blog) ([]bloggulus.Post, error)
	ReadPostBody(post bloggulus.Post) (string, error)
}

type reader struct{
	logger *log.Logger
}

func NewReader(logger *log.Logger) Reader {
	r := reader{
		logger: logger,
	}
	return &r
}

func (r *reader) ReadBlog(feedURL string) (bloggulus.Blog, error) {
	// early check to ensure the URL is valid
	_, err := url.Parse(feedURL)
	if err != nil {
		return bloggulus.Blog{}, err
	}

	// attempt to parse the feed via gofeed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return bloggulus.Blog{}, err
	}

	// create a bloggulus.Blog for the feed
	blog := bloggulus.NewBlog(feedURL, feed.Link, feed.Title)
	return blog, nil
}

func (r *reader) ReadBlogPosts(blog bloggulus.Blog) ([]bloggulus.Post, error) {
	// attempt to parse the feed via gofeed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(blog.FeedURL)
	if err != nil {
		return nil, err
	}

	// create a bloggulus.Post for each entry
	var posts []bloggulus.Post
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

		// ensure link is valid
		link := item.Link
		u, err := url.Parse(link)
		if err != nil {
			r.logger.Println(err)
			continue
		}

		// ensure link includes hostname
		if u.Hostname() == "" {
			link = feed.Link + link
		}

		// ensure link includes scheme
		matched, err := regexp.MatchString("^https?://", link)
		if err != nil {
			r.logger.Println(err)
			continue
		}
		// assume https if no scheme is present
		if !matched {
			link = "https://" + link
		}

		post := bloggulus.NewPost(link, item.Title, updated, blog)
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *reader) ReadPostBody(post bloggulus.Post) (string, error) {
	resp, err := http.Get(post.URL)
	if err != nil {
		return "", fmt.Errorf("%v: %v", post.URL, err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%v: %v", post.URL, err)
	}

	body := string(buf)
	body = CleanHTML(body)
	return body, nil
}
