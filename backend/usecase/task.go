package usecase

import (
	"backend/entity"
)

type taskUsecase struct {
	tr ITaskRepository
}

func NewTaskUsecase(tr ITaskRepository) ITaskUsecase {
	return &taskUsecase{tr: tr}
}

func (tu *taskUsecase) Create(task *entity.Task) (*entity.Task, error) {
	return tu.tr.Create(task)
}

func (tu *taskUsecase) Get(taskID entity.TaskID, userID entity.UserID) (*entity.Task, error) {
	return tu.tr.Get(taskID, userID)
}

func (tu *taskUsecase) GetAll(userID entity.UserID) (*[]entity.Task, error) {
	return tu.tr.GetAll(userID)
}

func (tu *taskUsecase) Save(task *entity.Task) (*entity.Task, error) {
	return tu.tr.Save(task)
}

func (tu *taskUsecase) Delete(taskID entity.TaskID, userID entity.UserID) error {
	return tu.tr.Delete(taskID, userID)
}
