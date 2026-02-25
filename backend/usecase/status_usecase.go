package usecase

import "backend/entity"

type IStatusUsecase interface {
	GetOrCreate(status *entity.Status) (*entity.Status, error)
}
