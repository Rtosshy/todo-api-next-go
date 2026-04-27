package handler

import (
	"backend/api"
	"backend/internal/infra/web/gin/presenter"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CsrfHandler interface {
	GetCsrfToken(c *gin.Context)
}

type csrfHandler struct{}

func NewCsrfHandler() CsrfHandler {
	return &csrfHandler{}
}

func (ch *csrfHandler) GetCsrfToken(c *gin.Context) {
	token := c.GetString("csrf")
	c.JSON(http.StatusOK, presenter.CsrfTokenResponse{
		ApiVersion: api.Version,
		Data:       presenter.CsrfToken(token),
	},
	)
}
