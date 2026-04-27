package entity_test

import (
	"backend/internal/domain"
	"backend/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	now := pkg.Str2time("2025-01-01")
	status := domain.Status{
		ID:   1,
		Name: "todo",
	}

	user := domain.User{
		ID:        1,
		Email:     "test@test.com",
		Password:  "password",
		CreatedAt: now,
	}

	task := domain.Task{
		ID:       1,
		Name:     "test",
		StatusID: 1,
		Status:   status,
		UserID:   1,
		User:     user,
	}
	assert.Equal(t, domain.TaskID(1), task.ID)
	assert.Equal(t, "test", task.Name)
	assert.Equal(t, domain.StatusID(1), task.StatusID)
	assert.Equal(t, domain.StatusID(1), task.Status.ID)
	assert.Equal(t, domain.StatusName("todo"), task.Status.Name)
	assert.Equal(t, domain.UserID(1), task.UserID)
	assert.Equal(t, domain.UserID(1), task.User.ID)
	assert.Equal(t, "test@test.com", task.User.Email)
	assert.Equal(t, "password", task.User.Password)
	assert.Equal(t, now, task.User.CreatedAt)
}
