package admin

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	id       uuid.UUID
	username string
	password string

	createdAt time.Time
	updatedAt time.Time
}

func NewAccount(username, password string) (*Account, error) {
	if username == "" {
		return nil, fmt.Errorf("account: invalid username")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	account := Account{
		id:       uuid.New(),
		username: username,
		password: string(hash),

		createdAt: now,
		updatedAt: now,
	}
	return &account, nil
}

func LoadAccount(id uuid.UUID, username, password string, createdAt, updatedAt time.Time) *Account {
	account := Account{
		id:       id,
		username: username,
		password: password,

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

func (a *Account) Password() string {
	return a.password
}

func (a *Account) PasswordMatches(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.password), []byte(password))
	return err == nil
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
