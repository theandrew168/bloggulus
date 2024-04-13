package mock

import (
	"github.com/theandrew168/bloggulus/backend/domain/reader"
	"github.com/theandrew168/bloggulus/backend/test"
)

func NewPost() *reader.Post {
	post := reader.LoadPost(
		test.RandomString(32),
		test.RandomURL(32),
		test.RandomString(32),
		test.RandomURL(32),
		test.RandomTime(),
		[]string{
			test.RandomString(8),
			test.RandomString(8),
			test.RandomString(8),
		},
	)
	return post
}
