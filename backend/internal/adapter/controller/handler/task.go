package handler

import (
	"backend/api"
	"backend/internal/adapter/controller/presenter"
	"backend/internal/domain"
	"backend/internal/usecase"
	"backend/pkg/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type taskHandler struct {
	tu usecase.TaskUsecase
}

func NewTaskHandler(tu usecase.TaskUsecase) TaskHandler {
	return &taskHandler{tu: tu}
}

func timeToDeadline(t *time.Time) *presenter.Deadline {
	if t == nil {
		return nil
	}
	deadline := presenter.Deadline{Time: *t}
	return &deadline
}

func deadlineToTime(d *presenter.Deadline) *time.Time {
	if d == nil {
		return nil
	}
	time := d.Time
	return &time
}

func taskToData(task *domain.Task) presenter.Task {
	return presenter.Task{
		Kind: "task",
		Id:   int(task.ID),
		Name: task.Name,
		Status: presenter.Status{
			Id:   (*int)(&task.Status.ID),
			Name: presenter.StatusName(task.Status.Name),
		},
		Deadline: timeToDeadline(task.Deadline),
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

func getUserIDFromContext(c *gin.Context) (domain.UserID, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		logger.Warn("user_id not found in context")
		return 0, fmt.Errorf("user_id not found in context")
	}

	userIDFloat, ok := userID.(float64)
	if !ok {
		logger.Warn(fmt.Sprintf("user_id has invalid type: %T", userID))
		return 0, fmt.Errorf("invalid user_id type")
	}

	return domain.UserID(userIDFloat), nil
}

func (th *taskHandler) CreateTask(c *gin.Context) {
	var requestBody presenter.CreateTaskRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	status, err := domain.NewStatus(string(requestBody.Status.Name))
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	task := &domain.Task{
		Name:     requestBody.Name,
		Status:   *status,
		UserID:   userID,
		Deadline: deadlineToTime(requestBody.Deadline),
	}

	createdTask, err := th.tu.Create(task)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, taskToResponse(createdTask))
}

func (th *taskHandler) GetTaskById(c *gin.Context, id int) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	taskID := domain.TaskID(id)

	task, err := th.tu.Get(taskID, userID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, taskToResponse(task))
}

func (th *taskHandler) GetAllTasks(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	tasks, err := th.tu.GetAll(userID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, tasksToResponse(tasks))
}

func (th *taskHandler) UpdateTaskById(c *gin.Context, id int) {
	var requestBody presenter.UpdateTaskRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	status, err := domain.NewStatus(string(requestBody.Status.Name))
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	taskID := domain.TaskID(id)

	task := &domain.Task{
		ID:       taskID,
		Name:     requestBody.Name,
		Status:   *status,
		UserID:   userID,
		Deadline: deadlineToTime(requestBody.Deadline),
	}

	updatedTask, err := th.tu.Save(task)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, taskToResponse(updatedTask))
}

func (th *taskHandler) DeleteTaskById(c *gin.Context, id int) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	taskID := domain.TaskID(id)

	if err := th.tu.Delete(taskID, userID); err != nil {
		c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
