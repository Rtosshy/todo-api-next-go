package usecase

import (
	"backend/internal/domain"
	"backend/internal/domain/repo"
	"backend/pkg/logger"
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type userUsecaseImpl struct {
	ur repo.UserRepo
}

func NewUserUsecase(ur repo.UserRepo) UserUsecase {
	return &userUsecaseImpl{ur: ur}
}

func (uu *userUsecaseImpl) SignUp(user *domain.User) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		logger.Error("Failed to hash password: " + err.Error())
		return nil, err
	}

	newUser := domain.User{
		Email:    user.Email,
		Password: string(hash),
	}

	return uu.ur.Create(&newUser)
}

func (uu *userUsecaseImpl) Login(user *domain.User) (string, error) {
	storedUser, err := uu.ur.GetByEmail(user.Email)
	if err != nil {
		logger.Error("GetByEmail failed: " + err.Error())
		return "", err
	}

	logger.Info(fmt.Sprintf("storedUser: ID=%d, Email=%s", storedUser.ID, storedUser.Email))

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		logger.Error("Password mismatch: " + err.Error())
		return "", err
	}

	logger.Info(fmt.Sprintf("Creating token with user_id: %d", storedUser.ID))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
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

func (uu *userUsecaseImpl) Save(user *domain.User) (*domain.User, error) {
	return uu.ur.Save(user)
}

func (uu *userUsecaseImpl) Delete(userID domain.UserID) error {
	return uu.ur.Delete(userID)
}
