package repo

import "backend/internal/domain"

type TaskRepo interface {
	Create(task *domain.Task) (*domain.Task, error)
	Get(taskID domain.TaskID, userID domain.UserID) (*domain.Task, error)
	GetAll(userID domain.UserID) (*[]domain.Task, error)
	Save(task *domain.Task) (*domain.Task, error)
	Delete(taskID domain.TaskID, userID domain.UserID) error
}
