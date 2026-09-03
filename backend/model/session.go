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

type NewSessionParams struct {
	Account *Account
	TTL     time.Duration
}

func NewSession(params NewSessionParams) (*Session, value.Token, error) {
	now := timeutil.Now()

	token, err := value.RandomToken()
	if err != nil {
		return nil, value.Token{}, err
	}

	session := Session{
		id:        uuid.New(),
		accountID: params.Account.ID(),
		tokenHash: token.Hash(),
		expiresAt: now.Add(params.TTL),
		meta:      NewMeta(),
	}

	return &session, token, nil
}

type LoadSessionParams struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	TokenHash value.TokenHash
	ExpiresAt time.Time
	Meta      *Meta
}

func LoadSession(params LoadSessionParams) *Session {
	session := Session{
		id:        params.ID,
		accountID: params.AccountID,
		tokenHash: params.TokenHash,
		expiresAt: params.ExpiresAt,
		meta:      params.Meta,
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
