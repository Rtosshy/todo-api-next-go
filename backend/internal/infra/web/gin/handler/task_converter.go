package handler

import (
	"time"

	"backend/api"
	"backend/internal/domain"
	"backend/internal/infra/web/gin/presenter"
	"backend/internal/usecase"
)

func deadlineToPresenter(d *domain.Deadline) *presenter.Deadline {
	if d == nil {
		return nil
	}
	return &presenter.Deadline{Time: d.Time()}
}

func presenterDeadlineToTimePtr(d *presenter.Deadline) *time.Time {
	if d == nil {
		return nil
	}
	t := d.Time
	return &t
}

func toCreateTaskInput(req presenter.CreateTaskRequestBody, userID domain.UserID) usecase.CreateTaskInput {
	return usecase.CreateTaskInput{
		Name:       req.Name,
		StatusName: string(req.Status.Name),
		UserID:     userID,
		Deadline:   presenterDeadlineToTimePtr(req.Deadline),
	}
}

func toUpdateTaskInput(req presenter.UpdateTaskRequestBody, taskID domain.TaskID, userID domain.UserID) usecase.UpdateTaskInput {
	return usecase.UpdateTaskInput{
		TaskID:     taskID,
		Name:       req.Name,
		StatusName: string(req.Status.Name),
		UserID:     userID,
		Deadline:   presenterDeadlineToTimePtr(req.Deadline),
	}
}

func taskToData(task *domain.Task) presenter.Task {
	statusID := int(task.Status().ID)
	return presenter.Task{
		Kind: "task",
		Id:   int(task.ID()),
		Name: task.Name().String(),
		Status: presenter.Status{
			Id:   &statusID,
			Name: presenter.StatusName(task.Status().Name.String()),
		},
		Deadline: deadlineToPresenter(task.Deadline()),
	}
}

func taskToResponse(task *domain.Task) presenter.TaskResponse {
	return presenter.TaskResponse{
		ApiVersion: api.Version,
		Data:       taskToData(task),
	}
}

func tasksToResponse(tasks *[]domain.Task) presenter.TasksResponse {
	data := make([]presenter.Task, len(*tasks))
	for i, task := range *tasks {
		data[i] = taskToData(&task)
	}
	return presenter.TasksResponse{
		ApiVersion: api.Version,
		Data:       data,
	}
}
