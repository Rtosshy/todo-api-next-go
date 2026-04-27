package db

import (
	"backend/internal/domain"
	"backend/internal/domain/repo"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type userRepoImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repo.UserRepo {
	return &userRepoImpl{db: db}
}

func (ur *userRepoImpl) Create(user *domain.User) (*domain.User, error) {
	if err := ur.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepoImpl) Get(userID domain.UserID) (*domain.User, error) {
	var user = domain.User{}
	if err := ur.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepoImpl) GetByEmail(email string) (*domain.User, error) {
	var user = domain.User{}
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepoImpl) Save(user *domain.User) (*domain.User, error) {
	selectedUser, err := ur.Get(user.ID)
	if err != nil {
		return nil, err
	}

	if err := copier.CopyWithOption(selectedUser, user, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}
	if err := ur.db.Save(selectedUser).Error; err != nil {
		return nil, err
	}

	return selectedUser, nil
}

func (ur *userRepoImpl) Delete(userID domain.UserID) error {
	user := domain.User{ID: userID}
	if err := ur.db.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}
