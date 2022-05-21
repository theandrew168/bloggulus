package main

import (
	"fmt"
	"os"

	"github.com/theandrew168/bloggulus"
)

func main() {
	os.Exit(run())
}

func run() int {
	client := bloggulus.NewClient(bloggulus.BaseURL)

	blogs, err := client.Blog.List()
	if err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Println("blogs:")
	for _, blog := range blogs {
		fmt.Println(blog.Title)
	}

	posts, err := client.Post.List()
	if err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Println("posts:")
	for _, post := range posts {
		fmt.Println(post.Title)
	}

	results, err := client.Post.Search("python")
	if err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Println("results:")
	for _, post := range results {
		fmt.Println(post.Title)
	}

	return 0
}
