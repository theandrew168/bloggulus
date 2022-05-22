package main

import (
	"log"

	"github.com/theandrew168/bloggulus"
)

func main() {
	client := bloggulus.NewClient(bloggulus.BaseURL)

	blogs, err := client.Blog.List()
	if err != nil {
		log.Fatalln(err)
	}

	for _, blog := range blogs {
		log.Println(blog.Title)
	}
}
