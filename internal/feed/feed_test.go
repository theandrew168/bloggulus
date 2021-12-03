package feed_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/internal/feed"
)

func TestCleanHTML(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"<code>hello world</code>", ""},
		{"<footer>hello world</footer>", ""},
		{"<header>hello world</header>", ""},
		{"<nav>hello world</nav>", ""},
		{"<pre>hello world</pre>", ""},

		{"hello world", "hello world"},
		{"<p>hello world</p>", "hello world"},
		{"<code><p>hello world</p></code>", ""},

		{"<script>console.log('hello')</script>", ""},
	}

	for _, test := range tests {
		if got := feed.CleanHTML(test.input); got != test.want {
			t.Errorf("feed.CleanHTML(%q) = %v", test.input, got)
		}
	}
}
