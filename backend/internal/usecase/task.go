package usecase

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
)

type taskUsecase struct {
	tr repository.TaskRepository
}

func NewTaskUsecase(tr repository.TaskRepository) *taskUsecase {
	return &taskUsecase{tr: tr}
}

func (tu *taskUsecase) Create(task *domain.Task) (*domain.Task, error) {
	return tu.tr.Create(task)
}

func (tu *taskUsecase) Get(taskID domain.TaskID, userID domain.UserID) (*domain.Task, error) {
	return tu.tr.Get(taskID, userID)
}

func (tu *taskUsecase) GetAll(userID domain.UserID) (*[]domain.Task, error) {
	return tu.tr.GetAll(userID)
}

func (tu *taskUsecase) Save(task *domain.Task) (*domain.Task, error) {
	return tu.tr.Save(task)
}

func (tu *taskUsecase) Delete(taskID domain.TaskID, userID domain.UserID) error {
	return tu.tr.Delete(taskID, userID)
}
