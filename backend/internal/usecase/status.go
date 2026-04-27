package usecase

import (
	"backend/internal/domain"
	"backend/internal/domain/repo"
)

type statusUsecase struct {
	sr repo.StatusRepo
}

func NewStatusUsecase(sr repo.StatusRepo) *statusUsecase {
	return &statusUsecase{sr: sr}
}

func (su *statusUsecase) GetOrCreate(status *domain.Status) (*domain.Status, error) {
	return su.sr.GetOrCreate(status)
}
