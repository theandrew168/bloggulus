package models

type Blog struct {
	BlogID   int
	FeedURL  string
	SiteURL  string
	Title    string

	Accounts []Account  // aka "Followers"
}
