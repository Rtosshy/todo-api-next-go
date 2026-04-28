package usecase

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"context"
)

type taskUsecase struct {
	tm TxManager
	tr repository.TaskRepository
}

func NewTaskUsecase(tm TxManager, tr repository.TaskRepository) *taskUsecase {
	return &taskUsecase{tm: tm, tr: tr}
}

func (tu *taskUsecase) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	return tu.tr.Create(ctx, task)
}

func (tu *taskUsecase) Get(ctx context.Context, taskID domain.TaskID, userID domain.UserID) (*domain.Task, error) {
	return tu.tr.Get(ctx, taskID, userID)
}

func (tu *taskUsecase) GetAll(ctx context.Context, userID domain.UserID) (*[]domain.Task, error) {
	return tu.tr.GetAll(ctx, userID)
}

func (tu *taskUsecase) Save(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	return tu.tr.Save(ctx, task)
}

func (tu *taskUsecase) Delete(ctx context.Context, taskID domain.TaskID, userID domain.UserID) error {
	return tu.tr.Delete(ctx, taskID, userID)
}
