package domain_test

import (
	"backend/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func validStatus(t *testing.T) domain.Status {
	t.Helper()
	s, err := domain.NewStatus("todo")
	assert.NoError(t, err)
	return s
}

func TestNewTask(t *testing.T) {
	taskName, err := domain.NewTaskName("Buy milk")
	assert.NoError(t, err)
	status := validStatus(t)

	t.Run("valid task without deadline", func(t *testing.T) {
		task, err := domain.NewTask(taskName, status, domain.UserID(1), nil)
		assert.NoError(t, err)
		assert.NotNil(t, task)
		assert.True(t, task.Name().Equals(taskName))
		assert.Equal(t, domain.UserID(1), task.UserID())
		assert.Nil(t, task.Deadline())
	})

	t.Run("valid task with deadline", func(t *testing.T) {
		dl, _ := domain.NewDeadline(time.Now().Add(24 * time.Hour))
		task, err := domain.NewTask(taskName, status, domain.UserID(1), &dl)
		assert.NoError(t, err)
		assert.NotNil(t, task.Deadline())
		assert.True(t, task.Deadline().Equals(dl))
	})

	t.Run("zero TaskName rejected", func(t *testing.T) {
		var zero domain.TaskName
		_, err := domain.NewTask(zero, status, domain.UserID(1), nil)
		assert.ErrorIs(t, err, domain.ErrInvalidTask)
	})

	t.Run("zero UserID rejected", func(t *testing.T) {
		_, err := domain.NewTask(taskName, status, domain.UserID(0), nil)
		assert.ErrorIs(t, err, domain.ErrInvalidTask)
	})
}

func TestReconstructTask(t *testing.T) {
	taskName, _ := domain.NewTaskName("Buy milk")
	status := validStatus(t)
	dl, _ := domain.NewDeadline(time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC))

	task := domain.ReconstructTask(domain.TaskID(7), taskName, status, domain.UserID(3), &dl)
	assert.Equal(t, domain.TaskID(7), task.ID())
	assert.Equal(t, "Buy milk", task.Name().String())
	assert.Equal(t, "todo", task.Status().Name.String())
	assert.Equal(t, domain.UserID(3), task.UserID())
	assert.NotNil(t, task.Deadline())
}
