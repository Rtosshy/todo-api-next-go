package handler

import "github.com/gin-gonic/gin"

type IUserHandler interface {
	PostSignUp(c *gin.Context)
	PostLogin(c *gin.Context)
	PostLogout(c *gin.Context)
}
