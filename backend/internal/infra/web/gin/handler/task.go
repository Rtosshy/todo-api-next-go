package handler

import (
	"backend/api"
	"backend/internal/domain"
	"backend/internal/infra/web/gin/presenter"
	"backend/pkg/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskUsecase interface {
	Create(task *domain.Task) (*domain.Task, error)
	Get(taskID domain.TaskID, userID domain.UserID) (*domain.Task, error)
	GetAll(userID domain.UserID) (*[]domain.Task, error)
	Save(task *domain.Task) (*domain.Task, error)
	Delete(taskID domain.TaskID, userID domain.UserID) error
}

type taskHandler struct {
	tu TaskUsecase
}

func NewTaskHandler(tu TaskUsecase) TaskHandler {
	return &taskHandler{tu: tu}
}

func deadlineToPresenter(d *domain.Deadline) *presenter.Deadline {
	if d == nil {
		return nil
	}
	return &presenter.Deadline{Time: d.Time()}
}

func presenterToDeadline(d *presenter.Deadline) (*domain.Deadline, error) {
	if d == nil {
		return nil, nil
	}
	dl, err := domain.NewDeadline(d.Time)
	if err != nil {
		return nil, err
	}
	return &dl, nil
}

func taskToData(task *domain.Task) presenter.Task {
	statusID := int(task.Status().ID)
	var deadlinePtr *time.Time
	if d := task.Deadline(); d != nil {
		t := d.Time()
		deadlinePtr = &t
	}
	var presenterDeadline *presenter.Deadline
	if deadlinePtr != nil {
		presenterDeadline = &presenter.Deadline{Time: *deadlinePtr}
	}
	return presenter.Task{
		Kind: "task",
		Id:   int(task.ID()),
		Name: task.Name().String(),
		Status: presenter.Status{
			Id:   &statusID,
			Name: presenter.StatusName(task.Status().Name.String()),
		},
		Deadline: presenterDeadline,
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

func buildTaskFromRequest(name string, statusName string, userID domain.UserID, deadline *presenter.Deadline) (*domain.Task, error) {
	taskName, err := domain.NewTaskName(name)
	if err != nil {
		return nil, err
	}
	status, err := domain.NewStatus(statusName)
	if err != nil {
		return nil, err
	}
	dl, err := presenterToDeadline(deadline)
	if err != nil {
		return nil, err
	}
	return domain.NewTask(taskName, status, userID, dl)
}

func (th *taskHandler) CreateTask(c *gin.Context) {
	var requestBody presenter.CreateTaskRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
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

	task, err := buildTaskFromRequest(requestBody.Name, string(requestBody.Status.Name), userID, requestBody.Deadline)
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
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

	task, err := th.tu.Get(domain.TaskID(id), userID)
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

	userID, err := getUserIDFromContext(c)
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	taskName, err := domain.NewTaskName(requestBody.Name)
	if err != nil {
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
	dl, err := presenterToDeadline(requestBody.Deadline)
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	task := domain.ReconstructTask(domain.TaskID(id), taskName, status, userID, dl)

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

	if err := th.tu.Delete(domain.TaskID(id), userID); err != nil {
		c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
