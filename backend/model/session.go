package model

import (
	"time"
	"uuid"

	"github.com/theandrew168/bloggulus/backend/timeutil"
	"github.com/theandrew168/bloggulus/backend/value"
)

type Session struct {
	id        uuid.UUID
	accountID uuid.UUID
	tokenHash value.TokenHash
	expiresAt time.Time

	meta *Meta
}

func NewSession(account *Account, ttl time.Duration) (*Session, value.Token, error) {
	now := timeutil.Now()

	token, err := value.RandomToken()
	if err != nil {
		return nil, value.Token{}, err
	}

	session := Session{
		id:        uuid.New(),
		accountID: account.ID(),
		tokenHash: token.Hash(),
		expiresAt: now.Add(ttl),

		meta: NewMeta(),
	}
	return &session, token, nil
}

func LoadSession(id, accountID uuid.UUID, tokenHash value.TokenHash, expiresAt time.Time, meta *Meta) *Session {
	session := Session{
		id:        id,
		accountID: accountID,
		tokenHash: tokenHash,
		expiresAt: expiresAt,

		meta: meta,
	}
	return &session
}

func (s *Session) ID() uuid.UUID {
	return s.id
}

func (s *Session) AccountID() uuid.UUID {
	return s.accountID
}

func (s *Session) TokenHash() value.TokenHash {
	return s.tokenHash
}

func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

func (s *Session) Meta() *Meta {
	return s.meta
}
