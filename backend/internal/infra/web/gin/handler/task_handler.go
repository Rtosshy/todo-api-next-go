package handler

import "github.com/gin-gonic/gin"

type TaskHandler interface {
	CreateTask(c *gin.Context)
	GetTaskById(c *gin.Context, id int)
	GetAllTasks(c *gin.Context)
	UpdateTaskById(c *gin.Context, id int)
	DeleteTaskById(c *gin.Context, id int)
}
