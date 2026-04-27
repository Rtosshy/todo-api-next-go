package domain

import (
	"errors"
	"net/mail"
)

var ErrInvalidEmail = errors.New("invalid email")

type Email struct {
	value string
}

func NewEmail(v string) (Email, error) {
	if v == "" {
		return Email{}, ErrInvalidEmail
	}
	addr, err := mail.ParseAddress(v)
	if err != nil {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: addr.Address}, nil
}

func (e Email) String() string { return e.value }

func (e Email) Equals(o Email) bool { return e.value == o.value }

func (e Email) IsZero() bool { return e.value == "" }
