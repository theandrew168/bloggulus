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

	for _, blog := range blogs {
		fmt.Println(blog.Title)
	}

	return 0
}
