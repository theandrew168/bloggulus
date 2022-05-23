package main

import (
	"log"

	"github.com/theandrew168/bloggulus"
)

const CustomURL = "https://bloggulus.com/api/v1"

func main() {
	client, err := bloggulus.NewClient(bloggulus.URL(CustomURL))
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
