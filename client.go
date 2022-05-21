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

type Client struct {
	Blog *BlogClient
	Post *PostClient
}

// NewClient returns a new Bloggulus API client.
func NewClient(url string) *Client {
	c := Client{
		Blog: NewBlogClient(url),
		Post: NewPostClient(url),
	}
	return &c
}

type BlogClient struct {
	client *http.Client
	url    string
}

// NewBlogClient returns a new Bloggulus blog client.
func NewBlogClient(url string) *BlogClient {
	c := BlogClient{
		client: new(http.Client),
		url:    url,
	}
	return &c
}

// Read reads a single blog by its ID.
func (c *BlogClient) Read(id int) (Blog, error) {
	endpoint := fmt.Sprintf("%s/blog/%d", c.url, id)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return Blog{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
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
	endpoint := fmt.Sprintf("%s/blog", c.url)
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

type PostClient struct {
	client *http.Client
	url    string
}

// NewPostClient returns a new Bloggulus post client.
func NewPostClient(url string) *PostClient {
	c := PostClient{
		client: new(http.Client),
		url:    url,
	}
	return &c
}

// Read reads a single post by its ID.
func (c *PostClient) Read(id int) (Post, error) {
	endpoint := fmt.Sprintf("%s/post/%d", c.url, id)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return Post{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
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
	endpoint := fmt.Sprintf("%s/post", c.url)
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
	endpoint := fmt.Sprintf("%s/post?q=%s", c.url, url.QueryEscape(query))
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
