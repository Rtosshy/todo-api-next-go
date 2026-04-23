package usecase

import (
	"backend/internal/domain"
	"backend/internal/domain/repo"
)

type taskUsecaseImpl struct {
	tr repo.TaskRepo
}

func NewTaskUsecase(tr repo.TaskRepo) TaskUsecase {
	return &taskUsecaseImpl{tr: tr}
}

func (tu *taskUsecaseImpl) Create(task *domain.Task) (*domain.Task, error) {
	return tu.tr.Create(task)
}

func (tu *taskUsecaseImpl) Get(taskID domain.TaskID, userID domain.UserID) (*domain.Task, error) {
	return tu.tr.Get(taskID, userID)
}

func (tu *taskUsecaseImpl) GetAll(userID domain.UserID) (*[]domain.Task, error) {
	return tu.tr.GetAll(userID)
}

func (tu *taskUsecaseImpl) Save(task *domain.Task) (*domain.Task, error) {
	return tu.tr.Save(task)
}

func (tu *taskUsecaseImpl) Delete(taskID domain.TaskID, userID domain.UserID) error {
	return tu.tr.Delete(taskID, userID)
}
