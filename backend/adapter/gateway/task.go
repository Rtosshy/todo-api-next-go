package gateway

import (
	"backend/entity"
	"backend/usecase"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) usecase.ITaskRepository {
	return &taskRepository{db: db}
}

func (tr *taskRepository) GetOrCreateStatus(task *entity.Task) error {
	var status entity.Status
	if err := tr.db.FirstOrCreate(&status, entity.Status{Name: task.Status.Name}).Error; err != nil {
		return err
	}
	task.StatusID = status.ID
	task.Status = status
	return nil
}

func (tr *taskRepository) Create(task *entity.Task) (*entity.Task, error) {
	if err := tr.GetOrCreateStatus(task); err != nil {
		return nil, err
	}
	if err := tr.db.Create(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (tr *taskRepository) Get(taskID entity.TaskID, userID entity.UserID) (*entity.Task, error) {
	var task = entity.Task{}
	if err := tr.db.Preload("Status").Preload("User").
		Where("id = ? AND user_id = ?", taskID, userID).
		First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (tr *taskRepository) GetAll(userID entity.UserID) (*[]entity.Task, error) {
	tasks := []entity.Task{}
	if err := tr.db.Preload("Status").Preload("User").
		Where("user_id = ?", userID).
		Order("created_at").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return &tasks, nil
}

func (tr *taskRepository) Save(task *entity.Task) (*entity.Task, error) {
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

func (tr *taskRepository) Delete(taskID entity.TaskID, userID entity.UserID) error {
	var task = entity.Task{}
	if err := tr.db.Where("id = ? AND user_id = ?", taskID, userID).Delete(&task).Error; err != nil {
		return err
	}
	return nil
}
