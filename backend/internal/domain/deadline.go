package domain

import (
	"errors"
	"time"
)

var ErrInvalidDeadline = errors.New("invalid deadline")

type Deadline struct {
	value time.Time
}

func NewDeadline(t time.Time) (Deadline, error) {
	if t.IsZero() {
		return Deadline{}, ErrInvalidDeadline
	}
	return Deadline{value: t}, nil
}

func (d Deadline) Time() time.Time { return d.value }

func (d Deadline) IsFuture() bool { return d.value.After(time.Now()) }

func (d Deadline) Equals(o Deadline) bool { return d.value.Equal(o.value) }
