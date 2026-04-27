package db

import (
	"backend/internal/domain"
	"backend/internal/domain/repo"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type taskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) repo.TaskRepo {
	return &taskRepo{db: db}
}

func (tr *taskRepo) GetOrCreateStatus(task *domain.Task) error {
	var status domain.Status
	if err := tr.db.FirstOrCreate(&status, domain.Status{Name: task.Status.Name}).Error; err != nil {
		return err
	}
	task.StatusID = status.ID
	task.Status = status
	return nil
}

func (tr *taskRepo) Create(task *domain.Task) (*domain.Task, error) {
	if err := tr.GetOrCreateStatus(task); err != nil {
		return nil, err
	}
	if err := tr.db.Create(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (tr *taskRepo) Get(taskID domain.TaskID, userID domain.UserID) (*domain.Task, error) {
	var task = domain.Task{}
	if err := tr.db.Preload("Status").Preload("User").
		Where("id = ? AND user_id = ?", taskID, userID).
		First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (tr *taskRepo) GetAll(userID domain.UserID) (*[]domain.Task, error) {
	tasks := []domain.Task{}
	if err := tr.db.Preload("Status").Preload("User").
		Where("user_id = ?", userID).
		Order("created_at").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return &tasks, nil
}

func (tr *taskRepo) Save(task *domain.Task) (*domain.Task, error) {
	selectedTask, err := tr.Get(task.ID, task.UserID)
	if err != nil {
		return nil, err
	}

	if err := tr.GetOrCreateStatus(task); err != nil {
		return nil, err
	}

	if err := copier.CopyWithOption(selectedTask, task, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}
	if err := tr.db.Save(selectedTask).Error; err != nil {
		return nil, err
	}

	return selectedTask, nil
}

func (tr *taskRepo) Delete(taskID domain.TaskID, userID domain.UserID) error {
	var task = domain.Task{}
	if err := tr.db.Where("id = ? AND user_id = ?", taskID, userID).Delete(&task).Error; err != nil {
		return err
	}
	return nil
}
