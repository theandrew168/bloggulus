package main

import (
	"log"

	"github.com/theandrew168/bloggulus"
)

func main() {
	client := bloggulus.NewClient(bloggulus.BaseURL)

	posts, err := client.Post.Search("python")
	if err != nil {
		log.Fatalln(err)
	}

	for _, post := range posts {
		log.Println(post.Title)
	}
}
