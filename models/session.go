package models

import (
	"time"
)

type Session struct {
	SessionID string
	AccountID string
	Expiry    time.Time
}
