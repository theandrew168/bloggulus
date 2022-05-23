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

	posts, err := client.Post.List()
	if err != nil {
		log.Fatalln(err)
	}

	for _, post := range posts {
		log.Println(post.Title)
	}
}
