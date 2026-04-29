package middleware

import (
	"net/http"

	"backend/internal/infra/web/gin/presenter"
	"backend/pkg/cookie"
	"backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CsrfTokenGenerator() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := uuid.New().String()
		c.Set("csrf", token)

		sameSite, secure, domain := cookie.GetCookieConfig()

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "_csrf",
			Value:    token,
			MaxAge:   3600,
			Path:     "/",
			Domain:   domain,
			Secure:   secure,
			HttpOnly: true,
			SameSite: sameSite,
		})

		c.Next()
	}
}

func CsrfValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStringCookie, err := c.Cookie("_csrf")
		if err != nil {
			logger.Warn("csrf token cookie not found: " + err.Error())
			c.JSON(presenter.NewErrorResponse(http.StatusForbidden, "csrf token required"))
			c.Abort()
			return
		}

		tokenStringHeader := c.GetHeader("X-CSRF-Token")
		if tokenStringHeader == "" {
			logger.Warn("csrf token header not found")
			c.JSON(presenter.NewErrorResponse(http.StatusForbidden, "csrf token required"))
			c.Abort()
			return
		}

		if tokenStringCookie != tokenStringHeader {
			logger.Warn("csrf token mismatch")
			c.JSON(presenter.NewErrorResponse(http.StatusForbidden, "invalid csrf token"))
			c.Abort()
			return
		}

		logger.Info("Csrf token validated!")

		c.Next()
	}
}
