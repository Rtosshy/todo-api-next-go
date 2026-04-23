package domain

import "time"

type TaskID int

type Task struct {
	ID        TaskID
	Name      string
	StatusID  StatusID
	Status    Status
	UserID    UserID
	User      User
	Deadline  *time.Time
	CreatedAt time.Time
}
