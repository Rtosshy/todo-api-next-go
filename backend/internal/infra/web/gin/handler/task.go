package handler

import (
	"backend/internal/domain"
	"backend/internal/infra/web/gin/presenter"
	"backend/internal/usecase"
	"backend/pkg/logger"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskUsecase interface {
	Create(ctx context.Context, in usecase.CreateTaskInput) (*domain.Task, error)
	Get(ctx context.Context, taskID domain.TaskID, userID domain.UserID) (*domain.Task, error)
	GetAll(ctx context.Context, userID domain.UserID) (*[]domain.Task, error)
	Save(ctx context.Context, in usecase.UpdateTaskInput) (*domain.Task, error)
	Delete(ctx context.Context, taskID domain.TaskID, userID domain.UserID) error
}

type taskHandler struct {
	tu TaskUsecase
}

func NewTaskHandler(tu TaskUsecase) TaskHandler {
	return &taskHandler{tu: tu}
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

	userID, err := getUserIDFromContext(c)
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	createdTask, err := th.tu.Create(c, toCreateTaskInput(requestBody, userID))
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

	task, err := th.tu.Get(c, domain.TaskID(id), userID)
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

	tasks, err := th.tu.GetAll(c, userID)
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

	updatedTask, err := th.tu.Save(c, toUpdateTaskInput(requestBody, domain.TaskID(id), userID))
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

	if err := th.tu.Delete(c, domain.TaskID(id), userID); err != nil {
		c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
