package middleware

import (
	"fmt"
	"net/http"
	"os"

	"backend/internal/adapter/controller/presenter"
	"backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			logger.Warn("Jwt token cookie not found: " + err.Error())
			c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, "authentication required"))
			c.Abort()
			return
		}

		logger.Info("Jwt token found, parsing...")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			logger.Warn("Jwt token parse error: " + err.Error())
			c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, "authentication failed"))
			c.Abort()
			return
		}

		if !token.Valid {
			logger.Warn("Jwt token is not valid")
			c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, "authentication failed"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Warn("Failed to parse claims")
			c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, "authentication failed"))
			c.Abort()
			return
		}

		logger.Info("Claims: " + fmt.Sprintf("%+v", claims))

		userID, exists := claims["user_id"]
		if !exists {
			logger.Warn("user_id not found in claims")
			c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, "authentication failed"))
			c.Abort()
			return
		}

		c.Set("user", token)
		c.Set("user_id", userID)
		logger.Info("user authenticated successfully with user_id: " + fmt.Sprintf("%v", userID))

		c.Next()
	}
}
