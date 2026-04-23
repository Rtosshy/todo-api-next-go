package usecase

import (
	"backend/internal/domain"
	"backend/internal/domain/repo"
)

type statusUsecaseImpl struct {
	sr repo.StatusRepo
}

func NewStatusUsecase(sr repo.StatusRepo) StatusUsecase {
	return &statusUsecaseImpl{sr: sr}
}

func (su statusUsecaseImpl) GetOrCreate(status *domain.Status) (*domain.Status, error) {
	return su.sr.GetOrCreate(status)
}
