package usecase

import "backend/entity"

type ITaskRepository interface {
	Create(task *entity.Task) (*entity.Task, error)
	Get(taskID entity.TaskID, userID entity.UserID) (*entity.Task, error)
	GetAll(userID entity.UserID) (*[]entity.Task, error)
	Save(task *entity.Task) (*entity.Task, error)
	Delete(taskID entity.TaskID, userID entity.UserID) error
}
