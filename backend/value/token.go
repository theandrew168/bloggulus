package value

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
)

var ErrEmptyToken = errors.New("token: value cannot be empty")
var ErrEmptyTokenHash = errors.New("token hash: value cannot be empty")

type Token struct {
	value string
}

func RandomToken() (Token, error) {
	t := Token{
		value: base64.RawURLEncoding.EncodeToString(randomBytes(32)),
	}
	return t, nil
}

func LoadToken(value string) (Token, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return Token{}, ErrEmptyToken
	}

	t := Token{
		value: trimmed,
	}
	return t, nil
}

func (t Token) Value() string {
	return t.value
}

func (t Token) Hash() TokenHash {
	hash := sha256.Sum256([]byte(t.value))
	th := TokenHash{
		value: hex.EncodeToString(hash[:]),
	}
	return th
}

type TokenHash struct {
	value string
}

func NewTokenHash(value string) (TokenHash, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return TokenHash{}, ErrEmptyTokenHash
	}

	t := TokenHash{
		value: trimmed,
	}
	return t, nil
}

func (t TokenHash) Value() string {
	return t.value
}

func randomBytes(n int) []byte {
	buf := make([]byte, n)

	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return buf
}
