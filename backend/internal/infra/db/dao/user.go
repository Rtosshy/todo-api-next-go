package dao

import "time"

type UserID int

type User struct {
	ID        UserID `gorm:"primaryKey"`
	Email     string `gorm:"unique"`
	Password  string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
