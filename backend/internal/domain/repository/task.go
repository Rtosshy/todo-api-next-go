package repository

import "backend/internal/domain/entity"

type TaskRepository interface {
	Create(task *entity.Task) (*entity.Task, error)
	Get(taskID entity.TaskID, userID entity.UserID) (*entity.Task, error)
	GetAll(userID entity.UserID) (*[]entity.Task, error)
	Save(task *entity.Task) (*entity.Task, error)
	Delete(taskID entity.TaskID, userID entity.UserID) error
}
