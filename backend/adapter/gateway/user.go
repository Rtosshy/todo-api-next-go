package gateway

import (
	"backend/entity"
	"backend/usecase"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) usecase.IUserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Create(user *entity.User) (*entity.User, error) {
	if err := ur.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) Get(userID entity.UserID) (*entity.User, error) {
	var user = entity.User{}
	if err := ur.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) GetByEmail(email string) (*entity.User, error) {
	var user = entity.User{}
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) Save(user *entity.User) (*entity.User, error) {
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

func (ur *userRepository) Delete(userID entity.UserID) error {
	user := entity.User{ID: userID}
	if err := ur.db.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}
