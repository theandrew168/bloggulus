package admin

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	id        uuid.UUID
	accountID uuid.UUID
	hash      string
	expiresAt time.Time

	createdAt time.Time
	updatedAt time.Time
}

func NewToken(account *Account, hash string, expiresAt time.Time) (*Token, error) {
	now := time.Now().UTC()
	if expiresAt.Before(now) {
		return nil, fmt.Errorf("token: expires in the past")
	}

	token := Token{
		id:        uuid.New(),
		accountID: account.ID(),
		hash:      hash,
		expiresAt: expiresAt,

		createdAt: now,
		updatedAt: now,
	}
	return &token, nil
}

func LoadToken(id, accountID uuid.UUID, hash string, expiresAt, createdAt, updatedAt time.Time) *Token {
	token := Token{
		id:        id,
		accountID: accountID,
		hash:      hash,
		expiresAt: expiresAt,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	return &token
}

func (t *Token) ID() uuid.UUID {
	return t.id
}

func (t *Token) AccountID() uuid.UUID {
	return t.accountID
}

func (t *Token) Hash() string {
	return t.hash
}

func (t *Token) ExpiresAt() time.Time {
	return t.expiresAt
}

func (t *Token) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Token) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Token) CheckDelete() error {
	return nil
}
