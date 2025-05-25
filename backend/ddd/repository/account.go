package repository

import "github.com/theandrew168/bloggulus/backend/ddd"

type Account interface {
	ReadBySessionID(sessionID string) (*ddd.Account, error)
}
