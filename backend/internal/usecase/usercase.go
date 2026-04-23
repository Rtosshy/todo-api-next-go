package usecase

import "backend/internal/domain"

type StatusUsecase interface {
	GetOrCreate(status *domain.Status) (*domain.Status, error)
}

type TaskUsecase interface {
	Create(task *domain.Task) (*domain.Task, error)
	Get(taskID domain.TaskID, userID domain.UserID) (*domain.Task, error)
	GetAll(userID domain.UserID) (*[]domain.Task, error)
	Save(task *domain.Task) (*domain.Task, error)
	Delete(taskID domain.TaskID, userID domain.UserID) error
}

type UserUsecase interface {
	SignUp(user *domain.User) (*domain.User, error)
	Login(user *domain.User) (string, error)
}
