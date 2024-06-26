package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/theandrew168/bloggulus/backend/timeutil"
)

type Account struct {
	id           uuid.UUID
	username     string
	passwordHash string
	isAdmin      bool

	createdAt time.Time
	updatedAt time.Time
}

func NewAccount(username, password string) (*Account, error) {
	if username == "" {
		return nil, fmt.Errorf("account: invalid username")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := timeutil.Now()
	account := Account{
		id:           uuid.New(),
		username:     username,
		passwordHash: string(passwordHash),
		isAdmin:      false,

		createdAt: now,
		updatedAt: now,
	}
	return &account, nil
}

func LoadAccount(id uuid.UUID, username, passwordHash string, isAdmin bool, createdAt, updatedAt time.Time) *Account {
	account := Account{
		id:           id,
		username:     username,
		passwordHash: passwordHash,
		isAdmin:      isAdmin,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	return &account
}

func (a *Account) ID() uuid.UUID {
	return a.id
}

func (a *Account) Username() string {
	return a.username
}

func (a *Account) PasswordHash() string {
	return a.passwordHash
}

func (a *Account) PasswordMatches(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.passwordHash), []byte(password))
	return err == nil
}

func (a *Account) IsAdmin() bool {
	return a.isAdmin
}

func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Account) UpdatedAt() time.Time {
	return a.updatedAt
}

func (a *Account) CheckDelete() error {
	return nil
}
