package usecase

import (
	"backend/entity"
	"backend/pkg/logger"
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	ur IUserRepository
}

func NewUserUsecase(ur IUserRepository) IUserUsecase {
	return &userUsecase{ur: ur}
}

func (uu *userUsecase) SignUp(user *entity.User) (*entity.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		logger.Error("Failed to hash password: " + err.Error())
		return nil, err
	}

	newUser := entity.User{
		Email:    user.Email,
		Password: string(hash),
	}

	return uu.ur.Create(&newUser)
}

func (uu *userUsecase) Login(user *entity.User) (string, error) {
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

func (uu *userUsecase) Save(user *entity.User) (*entity.User, error) {
	return uu.ur.Save(user)
}

func (uu *userUsecase) Delete(userID entity.UserID) error {
	return uu.ur.Delete(userID)
}
