package domain

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort      = errors.New("password must be at least 8 characters")
	ErrInvalidHashedPassword = errors.New("invalid hashed password")
)

const passwordMinLength = 8

type PlainPassword struct {
	value string
}

func NewPlainPassword(v string) (PlainPassword, error) {
	if len(v) < passwordMinLength {
		return PlainPassword{}, ErrPasswordTooShort
	}
	return PlainPassword{value: v}, nil
}

func (p PlainPassword) Hash() (HashedPassword, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.value), bcrypt.DefaultCost)
	if err != nil {
		return HashedPassword{}, err
	}
	return HashedPassword{value: string(hash)}, nil
}

func (p PlainPassword) Matches(h HashedPassword) bool {
	return bcrypt.CompareHashAndPassword([]byte(h.value), []byte(p.value)) == nil
}

type HashedPassword struct {
	value string
}

func NewHashedPassword(v string) (HashedPassword, error) {
	if v == "" {
		return HashedPassword{}, ErrInvalidHashedPassword
	}
	return HashedPassword{value: v}, nil
}

func (h HashedPassword) String() string { return h.value }

func (h HashedPassword) IsZero() bool { return h.value == "" }
