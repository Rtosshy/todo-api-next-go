package usecase

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"backend/pkg/logger"
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type userUsecase struct {
	ur repository.UserRepository
}

func NewUserUsecase(ur repository.UserRepository) *userUsecase {
	return &userUsecase{ur: ur}
}

func (uu *userUsecase) SignUp(email domain.Email, password domain.PlainPassword) (*domain.User, error) {
	hashed, err := password.Hash()
	if err != nil {
		logger.Error("Failed to hash password: " + err.Error())
		return nil, err
	}
	user, err := domain.NewUser(email, hashed)
	if err != nil {
		return nil, err
	}
	return uu.ur.Create(user)
}

func (uu *userUsecase) Login(email domain.Email, password domain.PlainPassword) (string, error) {
	storedUser, err := uu.ur.GetByEmail(email.String())
	if err != nil {
		logger.Error("GetByEmail failed: " + err.Error())
		return "", err
	}

	logger.Info(fmt.Sprintf("storedUser: ID=%d, Email=%s", storedUser.ID(), storedUser.Email().String()))

	if !password.Matches(storedUser.Password()) {
		logger.Error("Password mismatch")
		return "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID(),
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		logger.Error("Failed to sign token: " + err.Error())
		return "", err
	}

	logger.Info("Token created successfully")
	return tokenString, nil
}
