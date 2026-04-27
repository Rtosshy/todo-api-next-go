package domain

import (
	"errors"
	"unicode/utf8"
)

var ErrInvalidTaskName = errors.New("task name must be 1-255 characters")

const taskNameMaxLength = 255

type TaskName struct {
	value string
}

func NewTaskName(v string) (TaskName, error) {
	if v == "" {
		return TaskName{}, ErrInvalidTaskName
	}
	if utf8.RuneCountInString(v) > taskNameMaxLength {
		return TaskName{}, ErrInvalidTaskName
	}
	return TaskName{value: v}, nil
}

func (n TaskName) String() string { return n.value }

func (n TaskName) Equals(o TaskName) bool { return n.value == o.value }

func (n TaskName) IsZero() bool { return n.value == "" }
