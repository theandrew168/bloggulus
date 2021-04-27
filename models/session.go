package models

import (
	"time"
)

type Session struct {
	SessionID string
	AccountID int
	Expiry    time.Time

	Account   Account
}
