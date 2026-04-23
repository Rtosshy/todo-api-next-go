package domain

import "errors"

const (
	Todo       StatusName = "todo"
	InProgress StatusName = "inProgress"
	Done       StatusName = "done"
	Archive    StatusName = "archive"
	Pending    StatusName = "pending"
)

type StatusName string

func NewStatusName(value string) (*StatusName, error) {
	var statusName StatusName
	if err := statusName.Set(value); err != nil {
		return nil, err
	}
	return &statusName, nil
}

func (s *StatusName) IsValid() bool {
	return *s == Todo || *s == InProgress || *s == Done || *s == Archive || *s == Pending
}

func (s *StatusName) Set(value string) error {
	newStatusName := StatusName(value)
	if !newStatusName.IsValid() {
		return errors.New("Invalid value for StatusName")
	}
	*s = newStatusName
	return nil
}

type StatusID int

type Status struct {
	ID   StatusID
	Name StatusName
}

func NewStatus(name string) (*Status, error) {
	statusName, err := NewStatusName(name)
	if err != nil {
		return nil, err
	}
	return &Status{
		Name: *statusName,
	}, nil
}
