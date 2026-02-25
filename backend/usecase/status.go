package usecase

import (
	"backend/entity"
)

type statusUsecase struct {
	sr IStatusRepository
}

func NewStatusUsecase(sr IStatusRepository) IStatusUsecase {
	return &statusUsecase{sr: sr}
}

func (su statusUsecase) GetOrCreate(status *entity.Status) (*entity.Status, error) {
	return su.sr.GetOrCreateStatus(status)
}
