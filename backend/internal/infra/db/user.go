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

func userToDAO(u *domain.User) dao.User {
	return dao.User{
		ID:        dao.UserID(u.ID()),
		Email:     u.Email().String(),
		Password:  u.Password().String(),
		CreatedAt: u.CreatedAt(),
	}
}

func userToEntity(d *dao.User) (*domain.User, error) {
	email, err := domain.NewEmail(d.Email)
	if err != nil {
		return nil, err
	}
	hashed, err := domain.NewHashedPassword(d.Password)
	if err != nil {
		return nil, err
	}
	return domain.ReconstructUser(domain.UserID(d.ID), email, hashed, d.CreatedAt), nil
}

func (ur *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	d := userToDAO(user)
	if err := ur.db(ctx).Create(&d).Error; err != nil {
		return nil, err
	}
	return userToEntity(&d)
}

func (ur *userRepository) Get(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	var d dao.User
	if err := ur.db(ctx).First(&d, userID).Error; err != nil {
		return nil, err
	}
	return userToEntity(&d)
}

func (ur *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var d dao.User
	if err := ur.db(ctx).Where("email = ?", email).First(&d).Error; err != nil {
		return nil, err
	}
	return userToEntity(&d)
}

func (ur *userRepository) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
	d := userToDAO(user)
	if err := ur.db(ctx).Save(&d).Error; err != nil {
		return nil, err
	}
	return userToEntity(&d)
}

func (ur *userRepository) Delete(ctx context.Context, userID domain.UserID) error {
	return ur.db(ctx).Delete(&dao.User{ID: dao.UserID(userID)}).Error
}
