# bloggulus

[![Go Reference](https://pkg.go.dev/badge/github.com/theandrew168/bloggulus.svg)](https://pkg.go.dev/github.com/theandrew168/bloggulus)
[![Go Report Card](https://goreportcard.com/badge/github.com/theandrew168/bloggulus)](https://goreportcard.com/report/github.com/theandrew168/bloggulus)
[![MIT](https://img.shields.io/github/license/theandrew168/bloggulus)](https://img.shields.io/github/license/theandrew168/bloggulus)


Bloggulus is a web application for aggregating and indexing your favorite blogs.
I wrote it to serve as a less engaging and more personalized version of sites like Hacker News or Reddit.


## Client Library
This project includes a client library that can be used to programmatically read blogs and posts from a Bloggulus instance.
Here are a few basic example of how it works.


### List Blogs
```go
package main

import (
	"log"

	"github.com/theandrew168/bloggulus"
)

func main() {
	client, err := bloggulus.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	blogs, err := client.Blog.List()
	if err != nil {
		log.Fatalln(err)
	}

	for _, blog := range blogs {
		log.Println(blog.Title)
	}
}
```


### Search Posts
```go
package main

import (
	"log"

	"github.com/theandrew168/bloggulus"
)

func main() {
	client, err := bloggulus.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	posts, err := client.Post.Search("python")
	if err != nil {
		log.Fatalln(err)
	}

	for _, post := range posts {
		log.Println(post.Title)
	}
}
```


### Paginate Through Recent Posts
```go
package main

import (
	"log"

	"github.com/theandrew168/bloggulus"
)

func main() {
	client, err := bloggulus.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	// read 50 posts at a time starting at offset 0
	limit := 50
	offset := 0

	// paginate through all posts
	for {
		// read the current page
		posts, err := client.Post.List(
			client.Limit(limit),
			client.Offset(offset),
		)
		if err != nil {
			log.Fatalln(err)
		}

		// no more posts to read
		if len(posts) == 0 {
			break
		}

		// print current page of posts
		for _, post := range posts {
			log.Println(post.Title)
		}

		// update offset to prepare for next page
		offset += limit
	}
}
```


## Local Development
While the primary [Bloggulus website](https://bloggulus.com) represents my own personal collection of blogs, it is designed to be easily self-hostable.
Check out the [releases page](https://github.com/theandrew168/bloggulus/releases) for pre-built binaries and DEBs (for quick installation on Debian-based Linux systems)


### Setup
This project depends on the [Go programming language](https://golang.org/dl/) and the [TailwindCSS CLI](https://tailwindcss.com/blog/standalone-cli).


### Database
This project uses [PostgreSQL](https://www.postgresql.org/) for persistent storage.
To develop locally, you'll an instance of the database running somehow or another.
I find [Docker](https://www.docker.com/) to be a nice tool for this but you can do whatever works best.

The following command starts the necessary containers:
```bash
docker compose up -d
```

These containers can be stopped via:
```bash
docker compose down
```


### Running
If actively working on frontend templates, set `ENV=dev` to tell the server to reload templates from the filesystem on every page load.
Run the web server (in a background process) and let Tailwind watch for CSS changes:
```bash
# make web
ENV=dev go run cmd/web/main.go &
tailwindcss --watch -m -i tailwind.input.css -o internal/static/static/css/tailwind.min.css
```


### Testing
Tests can be ran after starting the necessary containers and applying database migrations:
```bash
# make test
go run cmd/web/main.go -migrate
go test -v ./...
```

Note that the tests will leave random test in the database so feel free to flush it out by restarting the containers:
```bash
docker compose down
docker compose up -d
```
