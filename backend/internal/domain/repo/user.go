package repo

import "backend/internal/domain"

type UserRepo interface {
	Create(user *domain.User) (*domain.User, error)
	Get(userID domain.UserID) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Save(user *domain.User) (*domain.User, error)
	Delete(userID domain.UserID) error
}
