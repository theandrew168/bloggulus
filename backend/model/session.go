package model

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"uuid"

	"github.com/theandrew168/bloggulus/backend/random"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

type Session struct {
	id        uuid.UUID
	accountID uuid.UUID
	tokenHash string
	expiresAt time.Time

	meta *Meta
}

// Generate a random, crypto-safe session token.
func GenerateSessionToken() (string, error) {
	return random.BytesBase64(32)
}

func NewSession(account *Account, ttl time.Duration) (*Session, string, error) {
	now := timeutil.Now()

	sessionToken, err := GenerateSessionToken()
	if err != nil {
		return nil, "", err
	}

	// Generate a SHA-256 hash of the plaintext session token. This will be the value
	// that we store in the `token_hash` field of our database table. Note that the
	// sha256.Sum256() function returns an array of length 32, so to make it easier to
	// work with we convert it to a slice using the [:] operator before storing it.
	sessionTokenHashBytes := sha256.Sum256([]byte(sessionToken))
	sessionTokenHash := hex.EncodeToString(sessionTokenHashBytes[:])

	session := Session{
		id:        uuid.New(),
		accountID: account.ID(),
		tokenHash: sessionTokenHash,
		expiresAt: now.Add(ttl),

		meta: NewMeta(),
	}
	return &session, sessionToken, nil
}

func LoadSession(id, accountID uuid.UUID, tokenHash string, expiresAt time.Time, meta *Meta) *Session {
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

func (s *Session) TokenHash() string {
	return s.tokenHash
}

func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

func (s *Session) Meta() *Meta {
	return s.meta
}
