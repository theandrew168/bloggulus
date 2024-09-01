package model

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

type Session struct {
	id        uuid.UUID
	accountID uuid.UUID
	hash      string
	expiresAt time.Time

	createdAt time.Time
	updatedAt time.Time
}

// Generate a random, crypto-safe session ID.
func GenerateSessionID() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func NewSession(account *Account, ttl time.Duration) (*Session, string, error) {
	now := timeutil.Now()

	sessionID, err := GenerateSessionID()
	if err != nil {
		return nil, "", err
	}

	// Generate a SHA-256 hash of the plaintext session ID. This will be the value
	// that we store in the `hash` field of our database table. Note that the
	// sha256.Sum256() function returns an array of length 32, so to make it easier to
	// work with we convert it to a slice using the [:] operator before storing it.
	hashBytes := sha256.Sum256([]byte(sessionID))
	hash := hex.EncodeToString(hashBytes[:])

	session := Session{
		id:        uuid.New(),
		accountID: account.ID(),
		hash:      hash,
		expiresAt: now.Add(ttl),

		createdAt: now,
		updatedAt: now,
	}
	return &session, sessionID, nil
}

func LoadSession(id, accountID uuid.UUID, hash string, expiresAt, createdAt, updatedAt time.Time) *Session {
	session := Session{
		id:        id,
		accountID: accountID,
		hash:      hash,
		expiresAt: expiresAt,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	return &session
}

func (s *Session) ID() uuid.UUID {
	return s.id
}

func (s *Session) AccountID() uuid.UUID {
	return s.accountID
}

func (s *Session) Hash() string {
	return s.hash
}

func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Session) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *Session) CheckDelete() error {
	return nil
}
