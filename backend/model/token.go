package model

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
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

func GenerateToken() (string, error) {
	// Initialize a zero-valued byte slice with a length of 16 bytes.
	randomBytes := make([]byte, 16)

	// Use the Read() function from the crypto/rand package to fill the byte slice with
	// random bytes from your operating system's CSPRNG. This will return an error if
	// the CSPRNG fails to function correctly.
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Note that by default base-32 strings may be padded at the end with the =
	// character. We don't need this padding character for the purpose of our tokens, so
	// we use the WithPadding(base32.NoPadding) method in the line below to omit them.
	token := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	return token, nil
}

func NewToken(account *Account, ttl time.Duration) (*Token, string, error) {
	now := time.Now().UTC().Round(time.Microsecond)

	value, err := GenerateToken()
	if err != nil {
		return nil, "", err
	}

	// Generate a SHA-256 hash of the plaintext token string. This will be the value
	// that we store in the `hash` field of our database table. Note that the
	// sha256.Sum256() function returns an *array* of length 32, so to make it easier to
	// work with we convert it to a slice using the [:] operator before storing it.
	hashBytes := sha256.Sum256([]byte(value))
	hash := hex.EncodeToString(hashBytes[:])

	token := Token{
		id:        uuid.New(),
		accountID: account.ID(),
		hash:      hash,
		expiresAt: now.Add(ttl),

		createdAt: now,
		updatedAt: now,
	}
	return &token, value, nil
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
