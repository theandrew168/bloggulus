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
