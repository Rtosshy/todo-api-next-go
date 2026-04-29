package usecase

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"context"
	"time"
)

type CreateTaskInput struct {
	Name       string
	StatusName string
	UserID     domain.UserID
	Deadline   *time.Time
}

type UpdateTaskInput struct {
	TaskID     domain.TaskID
	Name       string
	StatusName string
	UserID     domain.UserID
	Deadline   *time.Time
}

type taskUsecase struct {
	tm TxManager
	tr repository.TaskRepository
	sr repository.StatusRepository
}

func NewTaskUsecase(tm TxManager, tr repository.TaskRepository, sr repository.StatusRepository) *taskUsecase {
	return &taskUsecase{tm: tm, tr: tr, sr: sr}
}

func (tu *taskUsecase) Create(ctx context.Context, in CreateTaskInput) (*domain.Task, error) {
	var created *domain.Task
	err := tu.tm.RunInTx(ctx, func(ctx context.Context) error {
		taskName, err := domain.NewTaskName(in.Name)
		if err != nil {
			return err
		}
		status, err := domain.NewStatus(in.StatusName)
		if err != nil {
			return err
		}
		var deadline *domain.Deadline
		if in.Deadline != nil {
			dl, err := domain.NewDeadline(*in.Deadline)
			if err != nil {
				return err
			}
			deadline = &dl
		}

		resolvedStatus, err := tu.sr.GetOrCreate(ctx, &status)
		if err != nil {
			return err
		}

		task, err := domain.NewTask(taskName, *resolvedStatus, in.UserID, deadline)
		if err != nil {
			return err
		}

		created, err = tu.tr.Create(ctx, task)
		return err
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (tu *taskUsecase) Get(ctx context.Context, taskID domain.TaskID, userID domain.UserID) (*domain.Task, error) {
	return tu.tr.Get(ctx, taskID, userID)
}

func (tu *taskUsecase) GetAll(ctx context.Context, userID domain.UserID) (*[]domain.Task, error) {
	return tu.tr.GetAll(ctx, userID)
}

func (tu *taskUsecase) Save(ctx context.Context, in UpdateTaskInput) (*domain.Task, error) {
	var saved *domain.Task
	err := tu.tm.RunInTx(ctx, func(ctx context.Context) error {
		taskName, err := domain.NewTaskName(in.Name)
		if err != nil {
			return err
		}
		status, err := domain.NewStatus(in.StatusName)
		if err != nil {
			return err
		}
		var deadline *domain.Deadline
		if in.Deadline != nil {
			dl, err := domain.NewDeadline(*in.Deadline)
			if err != nil {
				return err
			}
			deadline = &dl
		}

		resolvedStatus, err := tu.sr.GetOrCreate(ctx, &status)
		if err != nil {
			return err
		}

		task := domain.ReconstructTask(in.TaskID, taskName, *resolvedStatus, in.UserID, deadline)

		saved, err = tu.tr.Save(ctx, task)
		return err
	})
	if err != nil {
		return nil, err
	}
	return saved, nil
}

func (tu *taskUsecase) Delete(ctx context.Context, taskID domain.TaskID, userID domain.UserID) error {
	return tu.tm.RunInTx(ctx, func(ctx context.Context) error {
		if _, err := tu.tr.Get(ctx, taskID, userID); err != nil {
			return err
		}
		return tu.tr.Delete(ctx, taskID, userID)
	})
}
