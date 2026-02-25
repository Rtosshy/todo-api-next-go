package usecase

import "backend/entity"

type IStatusRepository interface {
	GetOrCreateStatus(status *entity.Status) (*entity.Status, error)
}
