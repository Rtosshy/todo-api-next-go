package usecase

import "backend/entity"

type IUserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	Get(userID entity.UserID) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Save(user *entity.User) (*entity.User, error)
	Delete(userID entity.UserID) error
}
