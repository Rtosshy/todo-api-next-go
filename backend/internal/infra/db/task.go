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

func taskToDAO(task *domain.Task, statusID dao.StatusID) dao.Task {
	var deadline *time.Time
	if task.Deadline() != nil {
		v := task.Deadline().Time()
		deadline = &v
	}
	return dao.Task{
		ID:       dao.TaskID(task.ID()),
		Name:     task.Name().String(),
		StatusID: statusID,
		UserID:   dao.UserID(task.UserID()),
		Deadline: deadline,
	}
}

func taskToEntity(taskDAO *dao.Task) (*domain.Task, error) {
	name, err := domain.NewTaskName(taskDAO.Name)
	if err != nil {
		return nil, err
	}
	statusName, err := domain.NewStatusName(string(taskDAO.Status.Name))
	if err != nil {
		return nil, err
	}
	status := domain.Status{ID: domain.StatusID(taskDAO.Status.ID), Name: statusName}

	var deadline *domain.Deadline
	if taskDAO.Deadline != nil {
		dl, err := domain.NewDeadline(*taskDAO.Deadline)
		if err != nil {
			return nil, err
		}
		deadline = &dl
	}
	return domain.ReconstructTask(domain.TaskID(taskDAO.ID), name, status, domain.UserID(taskDAO.UserID), deadline), nil
}

func (tr *taskRepository) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	statusDAO := dao.Status{
		ID:   dao.StatusID(task.Status().ID),
		Name: dao.StatusName(task.Status().Name.String()),
	}
	taskDAO := taskToDAO(task, statusDAO.ID)
	if err := tr.db(ctx).Create(&taskDAO).Error; err != nil {
		return nil, err
	}
	taskDAO.Status = statusDAO
	return taskToEntity(&taskDAO)
}

func (tr *taskRepository) Get(ctx context.Context, taskID domain.TaskID, userID domain.UserID) (*domain.Task, error) {
	var taskDAO dao.Task
	if err := tr.db(ctx).Preload("Status").
		Where("id = ? AND user_id = ?", taskID, userID).
		First(&taskDAO).Error; err != nil {
		return nil, err
	}
	return taskToEntity(&taskDAO)
}

func (tr *taskRepository) GetAll(ctx context.Context, userID domain.UserID) (*[]domain.Task, error) {
	var taskDAOs []dao.Task
	if err := tr.db(ctx).Preload("Status").
		Where("user_id = ?", userID).
		Order("created_at").
		Find(&taskDAOs).Error; err != nil {
		return nil, err
	}
	tasks := make([]domain.Task, 0, len(taskDAOs))
	for i := range taskDAOs {
		task, err := taskToEntity(&taskDAOs[i])
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}
	return &tasks, nil
}

func (tr *taskRepository) Save(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	statusDAO := dao.Status{
		ID:   dao.StatusID(task.Status().ID),
		Name: dao.StatusName(task.Status().Name.String()),
	}
	taskDAO := taskToDAO(task, statusDAO.ID)
	if err := tr.db(ctx).Save(&taskDAO).Error; err != nil {
		return nil, err
	}
	taskDAO.Status = statusDAO
	return taskToEntity(&taskDAO)
}

func (tr *taskRepository) Delete(ctx context.Context, taskID domain.TaskID, userID domain.UserID) error {
	return tr.db(ctx).Where("id = ? AND user_id = ?", taskID, userID).Delete(&dao.Task{}).Error
}
