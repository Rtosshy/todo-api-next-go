package entity

import "time"

type UserID int

type User struct {
	id        UserID
	email     string
	password  string
	createdAt time.Time
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}
