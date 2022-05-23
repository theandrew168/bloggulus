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

	// offset posts by 5
	posts, err := client.Post.List(bloggulus.Offset(5))
	if err != nil {
		log.Fatalln(err)
	}

	for _, post := range posts {
		log.Println(post.Title)
	}
}
