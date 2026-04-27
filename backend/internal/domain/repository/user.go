package repository

import "backend/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	Get(userID entity.UserID) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Save(user *entity.User) (*entity.User, error)
	Delete(userID entity.UserID) error
}
