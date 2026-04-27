package handler

import (
	"backend/api"
	"backend/internal/domain"
	"backend/internal/infra/web/gin/presenter"
	"backend/pkg/cookie"
	"backend/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserUsecase interface {
	SignUp(user *domain.User) (*domain.User, error)
	Login(user *domain.User) (string, error)
}

type userHandler struct {
	uu UserUsecase
}

func NewUserHandler(uu UserUsecase) UserHandler {
	return &userHandler{uu: uu}
}

func (uh *userHandler) PostSignUp(c *gin.Context) {
	var requestBody presenter.SignUpRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	user := &domain.User{
		Email:    requestBody.User.Email,
		Password: *requestBody.User.Password,
	}

	// 平文パスワードを保存（Login用）
	plainPassword := *requestBody.User.Password

	createdUser, err := uh.uu.SignUp(user)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	// 自動ログインのためにLoginユースケースを呼び出す
	// 平文パスワードを持つUserオブジェクトを作成
	loginUser := &domain.User{
		Email:    createdUser.Email,
		Password: plainPassword,
	}
	tokenString, err := uh.uu.Login(loginUser)
	if err != nil {
		logger.Error((err.Error()))
		c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	sameSite, secure, domain := cookie.GetCookieConfig()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		MaxAge:   24 * 60 * 60,
		Path:     "/",
		Domain:   domain,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	})

	userID := int(createdUser.ID)

	c.JSON(http.StatusCreated, presenter.SignUpResponse{
		ApiVersion: api.Version,
		Data: presenter.User{
			Kind:  "user",
			Id:    &userID,
			Email: createdUser.Email,
		},
	})
}

func (uh *userHandler) PostLogin(c *gin.Context) {
	var requestBody presenter.LoginRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	user := &domain.User{
		Email:    requestBody.User.Email,
		Password: *requestBody.User.Password,
	}

	tokenString, err := uh.uu.Login(user)
	if err != nil {
		logger.Error((err.Error()))
		c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	sameSite, secure, domain := cookie.GetCookieConfig()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		MaxAge:   24 * 60 * 60,
		Path:     "/",
		Domain:   domain,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	})
	c.Status(http.StatusOK)
}

func (uh *userHandler) PostLogout(c *gin.Context) {
	sameSite, secure, domain := cookie.GetCookieConfig()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   domain,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	})
	c.Status(http.StatusOK)
}
