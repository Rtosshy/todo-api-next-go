package db

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"backend/internal/infra/db/dao"
	"time"

	"context"
)

type taskRepository struct {
	baseRepository
}

func NewTaskRepository(b baseRepository) repository.TaskRepository {
	return &taskRepository{b}
}

func (tr *taskRepository) getOrCreateStatus(ctx context.Context, name dao.StatusName) (dao.Status, error) {
	var s dao.Status
	if err := tr.db(ctx).FirstOrCreate(&s, dao.Status{Name: name}).Error; err != nil {
		return dao.Status{}, err
	}
	return s, nil
}

func taskToDAO(t *domain.Task, statusID dao.StatusID) dao.Task {
	var deadline *time.Time
	if t.Deadline() != nil {
		v := t.Deadline().Time()
		deadline = &v
	}
	return dao.Task{
		ID:       dao.TaskID(t.ID()),
		Name:     t.Name().String(),
		StatusID: statusID,
		UserID:   dao.UserID(t.UserID()),
		Deadline: deadline,
	}
}

func taskToEntity(d *dao.Task) (*domain.Task, error) {
	name, err := domain.NewTaskName(d.Name)
	if err != nil {
		return nil, err
	}
	statusName, err := domain.NewStatusName(string(d.Status.Name))
	if err != nil {
		return nil, err
	}
	status := domain.Status{ID: domain.StatusID(d.Status.ID), Name: statusName}

	var deadline *domain.Deadline
	if d.Deadline != nil {
		dl, err := domain.NewDeadline(*d.Deadline)
		if err != nil {
			return nil, err
		}
		deadline = &dl
	}
	return domain.ReconstructTask(domain.TaskID(d.ID), name, status, domain.UserID(d.UserID), deadline), nil
}

func (tr *taskRepository) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	statusDAO, err := tr.getOrCreateStatus(ctx, dao.StatusName(task.Status().Name.String()))
	if err != nil {
		return nil, err
	}
	d := taskToDAO(task, statusDAO.ID)
	if err := tr.db(ctx).Create(&d).Error; err != nil {
		return nil, err
	}
	d.Status = statusDAO
	return taskToEntity(&d)
}

func (tr *taskRepository) Get(ctx context.Context, taskID domain.TaskID, userID domain.UserID) (*domain.Task, error) {
	var d dao.Task
	if err := tr.db(ctx).Preload("Status").
		Where("id = ? AND user_id = ?", taskID, userID).
		First(&d).Error; err != nil {
		return nil, err
	}
	return taskToEntity(&d)
}

func (tr *taskRepository) GetAll(ctx context.Context, userID domain.UserID) (*[]domain.Task, error) {
	var ds []dao.Task
	if err := tr.db(ctx).Preload("Status").
		Where("user_id = ?", userID).
		Order("created_at").
		Find(&ds).Error; err != nil {
		return nil, err
	}
	tasks := make([]domain.Task, 0, len(ds))
	for i := range ds {
		t, err := taskToEntity(&ds[i])
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}
	return &tasks, nil
}

func (tr *taskRepository) Save(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	statusDAO, err := tr.getOrCreateStatus(ctx, dao.StatusName(task.Status().Name.String()))
	if err != nil {
		return nil, err
	}
	d := taskToDAO(task, statusDAO.ID)
	if err := tr.db(ctx).Save(&d).Error; err != nil {
		return nil, err
	}
	d.Status = statusDAO
	return taskToEntity(&d)
}

func (tr *taskRepository) Delete(ctx context.Context, taskID domain.TaskID, userID domain.UserID) error {
	return tr.db(ctx).Where("id = ? AND user_id = ?", taskID, userID).Delete(&dao.Task{}).Error
}
