// Package bloggulus provides a client library for a deployed Bloggulus
// server (either bloggulus.com or a self-hosted instance). The library
// communicates with the server via its REST API. No changes can be made
// to the server with this library: it only allows you to read data out.
package bloggulus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// BaseURL points to the publicly accessible Bloggulus API.
	BaseURL = "https://bloggulus.com/api/v1"
)

// Client holds the state necessary to communicate with a Bloggulus API.
type Client struct {
	client  *http.Client
	baseURL string

	Blog *BlogClient
	Post *PostClient
}

// ClientOption represents a functional option of the Client type.
type ClientOption func(*Client) error

// NewClient returns a new Bloggulus API client.
func NewClient(options ...ClientOption) (*Client, error) {
	c := Client{
		client:  new(http.Client),
		baseURL: BaseURL,
	}
	c.Blog = NewBlogClient(c)
	c.Post = NewPostClient(c)

	for _, option := range options {
		err := option(&c)
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// URL sets the base URL for this client.
func URL(url string) ClientOption {
	return func(c *Client) error {
		c.baseURL = url
		return nil
	}
}

func (c *Client) get(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// BlogClient is a Client for Bloggulus blogs.
type BlogClient struct {
	Client
}

// NewBlogClient returns a new Bloggulus blog client.
func NewBlogClient(client Client) *BlogClient {
	c := BlogClient{
		client,
	}
	return &c
}

// Read reads a single blog by its ID.
func (c *BlogClient) Read(id int) (Blog, error) {
	endpoint := fmt.Sprintf("%s/blog/%d", c.baseURL, id)
	resp, err := c.get(endpoint)
	if err != nil {
		return Blog{}, err
	}
	defer resp.Body.Close()

	var msg struct {
		Blog Blog `json:"blog"`
	}
	err = json.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		return Blog{}, err
	}

	return msg.Blog, nil
}

// List lists all blogs in alphabetical order by title.
func (c *BlogClient) List() ([]Blog, error) {
	endpoint := fmt.Sprintf("%s/blog", c.baseURL)
	resp, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var msg struct {
		Blogs []Blog `json:"blogs"`
	}
	err = json.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		return nil, err
	}

	return msg.Blogs, nil
}

// PostClient is a Client for Bloggulus posts.
type PostClient struct {
	Client
}

// NewPostClient returns a new Bloggulus post client.
func NewPostClient(client Client) *PostClient {
	c := PostClient{
		client,
	}
	return &c
}

// Read reads a single post by its ID.
func (c *PostClient) Read(id int) (Post, error) {
	endpoint := fmt.Sprintf("%s/post/%d", c.baseURL, id)
	resp, err := c.get(endpoint)
	if err != nil {
		return Post{}, err
	}
	defer resp.Body.Close()

	var msg struct {
		Post Post `json:"post"`
	}
	err = json.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		return Post{}, err
	}

	return msg.Post, nil
}

// List lists all posts in reverse chronological orders (newest first).
func (c *PostClient) List() ([]Post, error) {
	endpoint := fmt.Sprintf("%s/post", c.baseURL)
	resp, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var msg struct {
		Posts []Post `json:"posts"`
	}
	err = json.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		return nil, err
	}

	return msg.Posts, nil
}

// Search searches all posts based on a given query string.
func (c *PostClient) Search(query string) ([]Post, error) {
	endpoint := fmt.Sprintf("%s/post?q=%s", c.baseURL, url.QueryEscape(query))
	resp, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var msg struct {
		Posts []Post `json:"posts"`
	}
	err = json.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		return nil, err
	}

	return msg.Posts, nil
}
