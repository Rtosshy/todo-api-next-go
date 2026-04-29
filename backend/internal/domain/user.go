package domain

import (
	"errors"
	"time"
)

var ErrInvalidUser = errors.New("invalid user")

type UserID int

type User struct {
	id        UserID
	email     Email
	password  HashedPassword
	createdAt time.Time
}

func NewUser(email Email, password HashedPassword) (*User, error) {
	if email.IsZero() || password.IsZero() {
		return nil, ErrInvalidUser
	}
	return &User{email: email, password: password}, nil
}

func ReconstructUser(id UserID, email Email, password HashedPassword, createdAt time.Time) *User {
	return &User{id: id, email: email, password: password, createdAt: createdAt}
}

func (u *User) ID() UserID              { return u.id }
func (u *User) Email() Email            { return u.email }
func (u *User) Password() HashedPassword { return u.password }
func (u *User) CreatedAt() time.Time    { return u.createdAt }
