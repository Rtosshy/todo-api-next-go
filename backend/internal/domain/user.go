package domain

import "time"

type UserID int

type User struct {
	ID        UserID
	Email     string
	Password  string
	CreatedAt time.Time
}
