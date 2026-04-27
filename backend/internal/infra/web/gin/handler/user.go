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
	SignUp(email domain.Email, password domain.PlainPassword) (*domain.User, error)
	Login(email domain.Email, password domain.PlainPassword) (string, error)
}

type userHandler struct {
	uu UserUsecase
}

func NewUserHandler(uu UserUsecase) UserHandler {
	return &userHandler{uu: uu}
}

func bindCredentials(emailStr string, passwordStr *string) (domain.Email, domain.PlainPassword, error) {
	email, err := domain.NewEmail(emailStr)
	if err != nil {
		return domain.Email{}, domain.PlainPassword{}, err
	}
	pwStr := ""
	if passwordStr != nil {
		pwStr = *passwordStr
	}
	password, err := domain.NewPlainPassword(pwStr)
	if err != nil {
		return domain.Email{}, domain.PlainPassword{}, err
	}
	return email, password, nil
}

func (uh *userHandler) PostSignUp(c *gin.Context) {
	var requestBody presenter.SignUpRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	email, password, err := bindCredentials(requestBody.User.Email, requestBody.User.Password)
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	createdUser, err := uh.uu.SignUp(email, password)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	tokenString, err := uh.uu.Login(email, password)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	sameSite, secure, cookieDomain := cookie.GetCookieConfig()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		MaxAge:   24 * 60 * 60,
		Path:     "/",
		Domain:   cookieDomain,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	})

	userID := int(createdUser.ID())

	c.JSON(http.StatusCreated, presenter.SignUpResponse{
		ApiVersion: api.Version,
		Data: presenter.User{
			Kind:  "user",
			Id:    &userID,
			Email: createdUser.Email().String(),
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

	email, password, err := bindCredentials(requestBody.User.Email, requestBody.User.Password)
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	tokenString, err := uh.uu.Login(email, password)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	sameSite, secure, cookieDomain := cookie.GetCookieConfig()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		MaxAge:   24 * 60 * 60,
		Path:     "/",
		Domain:   cookieDomain,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	})
	c.Status(http.StatusOK)
}

func (uh *userHandler) PostLogout(c *gin.Context) {
	sameSite, secure, cookieDomain := cookie.GetCookieConfig()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   cookieDomain,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	})
	c.Status(http.StatusOK)
}
