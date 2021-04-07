package models

type Account struct {
	AccountID int
	Username  string
	Password  string
	Email     string
	Verified  bool

	Blogs     []Blog  // aka "Follows"
	Sessions  []Session
}
