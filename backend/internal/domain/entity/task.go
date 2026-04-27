package entity

import "time"

type TaskID int

type Task struct {
	id        TaskID
	name      string
	statusID  StatusID
	status    Status
	userID    UserID
	user      User
	deadline  *time.Time
	createdAt time.Time
}

func (t *Task) ID() TaskID {
	return t.id
}
func (t *Task) Name() string {
	return t.name
}
func (t *Task) StatusID() StatusID {
	return t.statusID
}
func (t *Task) Status() Status {
	return t.status
}
func (t *Task) UserID() UserID {
	return t.userID
}
func (t *Task) User() User {
	return t.user
}
func (t *Task) Deadline() *time.Time {
	return t.deadline
}
func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}
