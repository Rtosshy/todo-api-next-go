package usecase

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
)

type statusUsecase struct {
	sr repository.StatusRepository
}

func NewStatusUsecase(sr repository.StatusRepository) *statusUsecase {
	return &statusUsecase{sr: sr}
}

func (su *statusUsecase) GetOrCreate(status *domain.Status) (*domain.Status, error) {
	return su.sr.GetOrCreate(status)
}
