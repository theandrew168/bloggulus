package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/timeutil"
)

type Account struct {
	id      uuid.UUID
	email   string
	isAdmin bool

	createdAt time.Time
	updatedAt time.Time
}

func NewAccount(email string) (*Account, error) {
	if email == "" {
		return nil, fmt.Errorf("account: invalid email")
	}

	now := timeutil.Now()
	account := Account{
		id:      uuid.New(),
		email:   email,
		isAdmin: false,

		createdAt: now,
		updatedAt: now,
	}
	return &account, nil
}

func LoadAccount(id uuid.UUID, email string, isAdmin bool, createdAt, updatedAt time.Time) *Account {
	account := Account{
		id:      id,
		email:   email,
		isAdmin: isAdmin,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	return &account
}

func (a *Account) ID() uuid.UUID {
	return a.id
}

func (a *Account) Email() string {
	return a.email
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
