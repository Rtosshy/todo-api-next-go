package repository

import (
	"context"

	"backend/internal/domain"
)

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) (*domain.Task, error)
	Get(ctx context.Context, taskID domain.TaskID, userID domain.UserID) (*domain.Task, error)
	GetAll(ctx context.Context, userID domain.UserID) (*[]domain.Task, error)
	Save(ctx context.Context, task *domain.Task) (*domain.Task, error)
	Delete(ctx context.Context, taskID domain.TaskID, userID domain.UserID) error
}
