package db

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"backend/internal/infra/db/dao"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
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

func (ur *userRepository) Create(user *domain.User) (*domain.User, error) {
	d := userToDAO(user)
	if err := ur.db.Create(&d).Error; err != nil {
		return nil, err
	}
	return userToEntity(&d)
}

func (ur *userRepository) Get(userID domain.UserID) (*domain.User, error) {
	var d dao.User
	if err := ur.db.First(&d, userID).Error; err != nil {
		return nil, err
	}
	return userToEntity(&d)
}

func (ur *userRepository) GetByEmail(email string) (*domain.User, error) {
	var d dao.User
	if err := ur.db.Where("email = ?", email).First(&d).Error; err != nil {
		return nil, err
	}
	return userToEntity(&d)
}

func (ur *userRepository) Save(user *domain.User) (*domain.User, error) {
	d := userToDAO(user)
	if err := ur.db.Save(&d).Error; err != nil {
		return nil, err
	}
	return userToEntity(&d)
}

func (ur *userRepository) Delete(userID domain.UserID) error {
	return ur.db.Delete(&dao.User{ID: dao.UserID(userID)}).Error
}
