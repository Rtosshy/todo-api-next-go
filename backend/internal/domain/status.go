package domain

import "errors"

var ErrInvalidStatusName = errors.New("invalid status name")

const (
	statusTodo       = "todo"
	statusInProgress = "inProgress"
	statusDone       = "done"
	statusArchive    = "archive"
	statusPending    = "pending"
)

type StatusName struct {
	value string
}

func NewStatusName(v string) (StatusName, error) {
	switch v {
	case statusTodo, statusInProgress, statusDone, statusArchive, statusPending:
		return StatusName{value: v}, nil
	}
	return StatusName{}, ErrInvalidStatusName
}

func (s StatusName) String() string { return s.value }

func (s StatusName) Equals(o StatusName) bool { return s.value == o.value }

func (s StatusName) IsZero() bool { return s.value == "" }

var (
	StatusNameTodo       = StatusName{value: statusTodo}
	StatusNameInProgress = StatusName{value: statusInProgress}
	StatusNameDone       = StatusName{value: statusDone}
	StatusNameArchive    = StatusName{value: statusArchive}
	StatusNamePending    = StatusName{value: statusPending}
)

type StatusID int

type Status struct {
	ID   StatusID
	Name StatusName
}

func NewStatus(name string) (Status, error) {
	statusName, err := NewStatusName(name)
	if err != nil {
		return Status{}, err
	}
	return Status{Name: statusName}, nil
}
