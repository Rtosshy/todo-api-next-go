package usecase

import "backend/entity"

type IUserUsecase interface {
	SignUp(user *entity.User) (*entity.User, error)
	Login(user *entity.User) (string, error)
}
