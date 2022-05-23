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
			bloggulus.Limit(limit),
			bloggulus.Offset(offset),
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
