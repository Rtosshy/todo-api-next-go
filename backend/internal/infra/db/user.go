package db

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"backend/internal/infra/db/dao"

	"context"
)

type userRepository struct {
	baseRepository
}

func NewUserRepository(b baseRepository) repository.UserRepository {
	return &userRepository{b}
}

func userToDAO(user *domain.User) dao.User {
	return dao.User{
		ID:        dao.UserID(user.ID()),
		Email:     user.Email().String(),
		Password:  user.Password().String(),
		CreatedAt: user.CreatedAt(),
	}
}

func userToEntity(userDAO *dao.User) (*domain.User, error) {
	email, err := domain.NewEmail(userDAO.Email)
	if err != nil {
		return nil, err
	}
	hashed, err := domain.NewHashedPassword(userDAO.Password)
	if err != nil {
		return nil, err
	}
	return domain.ReconstructUser(domain.UserID(userDAO.ID), email, hashed, userDAO.CreatedAt), nil
}

func (ur *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	userDAO := userToDAO(user)
	if err := ur.db(ctx).Create(&userDAO).Error; err != nil {
		return nil, err
	}
	return userToEntity(&userDAO)
}

func (ur *userRepository) Get(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	var userDAO dao.User
	if err := ur.db(ctx).First(&userDAO, userID).Error; err != nil {
		return nil, err
	}
	return userToEntity(&userDAO)
}

func (ur *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var userDAO dao.User
	if err := ur.db(ctx).Where("email = ?", email).First(&userDAO).Error; err != nil {
		return nil, err
	}
	return userToEntity(&userDAO)
}

func (ur *userRepository) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
	userDAO := userToDAO(user)
	if err := ur.db(ctx).Save(&userDAO).Error; err != nil {
		return nil, err
	}
	return userToEntity(&userDAO)
}

func (ur *userRepository) Delete(ctx context.Context, userID domain.UserID) error {
	return ur.db(ctx).Delete(&dao.User{ID: dao.UserID(userID)}).Error
}
