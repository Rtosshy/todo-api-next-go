package domain

import "errors"

var ErrInvalidTask = errors.New("invalid task")

type TaskID int

type Task struct {
	id       TaskID
	name     TaskName
	status   Status
	userID   UserID
	deadline *Deadline
}

func NewTask(name TaskName, status Status, userID UserID, deadline *Deadline) (*Task, error) {
	if name.IsZero() || userID == 0 {
		return nil, ErrInvalidTask
	}
	return &Task{name: name, status: status, userID: userID, deadline: deadline}, nil
}

func ReconstructTask(id TaskID, name TaskName, status Status, userID UserID, deadline *Deadline) *Task {
	return &Task{id: id, name: name, status: status, userID: userID, deadline: deadline}
}

func (t *Task) ID() TaskID          { return t.id }
func (t *Task) Name() TaskName      { return t.name }
func (t *Task) Status() Status      { return t.status }
func (t *Task) UserID() UserID      { return t.userID }
func (t *Task) Deadline() *Deadline { return t.deadline }
