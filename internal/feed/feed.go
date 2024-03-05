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

	"github.com/theandrew168/bloggulus/internal/domain"
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
	ReadBlog(feedURL string) (domain.Blog, error)
	ReadBlogPosts(blog domain.Blog, body io.Reader) ([]domain.Post, error)
	ReadPostBody(post domain.Post) (string, error)
}

type reader struct {
	logger *log.Logger
}

func NewReader(logger *log.Logger) Reader {
	r := reader{
		logger: logger,
	}
	return &r
}

func (r *reader) ReadBlog(feedURL string) (domain.Blog, error) {
	// early check to ensure the URL is valid
	_, err := url.Parse(feedURL)
	if err != nil {
		return domain.Blog{}, err
	}

	resp, err := http.Get(feedURL)
	if err != nil {
		return domain.Blog{}, err
	}
	defer resp.Body.Close()

	etag := resp.Header.Get("ETag")
	lastModified := resp.Header.Get("Last-Modified")

	// attempt to parse the feed via gofeed
	fp := gofeed.NewParser()
	feed, err := fp.Parse(resp.Body)
	if err != nil {
		return domain.Blog{}, err
	}

	// create a domain.Blog for the feed
	blog := domain.NewBlog(feedURL, feed.Link, feed.Title, etag, lastModified)
	return blog, nil
}

func (r *reader) ReadBlogPosts(blog domain.Blog, body io.Reader) ([]domain.Post, error) {
	// attempt to parse the feed via gofeed
	fp := gofeed.NewParser()
	feed, err := fp.Parse(body)
	if err != nil {
		return nil, err
	}

	// create a domain.Post for each entry
	var posts []domain.Post
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

		post := domain.NewPost(link, item.Title, updated, item.Content, blog)
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *reader) ReadPostBody(post domain.Post) (string, error) {
	// fetch post body if it wasn't included in the feed
	body := post.Body
	if body == "" {
		resp, err := http.Get(post.URL)
		if err != nil {
			return "", fmt.Errorf("%v: %v", post.URL, err)
		}
		defer resp.Body.Close()

		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("%v: %v", post.URL, err)
		}

		body = string(buf)
	}

	body = CleanHTML(body)
	return body, nil
}
