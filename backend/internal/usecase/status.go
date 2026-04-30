package usecase

import (
	"context"

	"backend/internal/domain"
	"backend/internal/domain/repository"
)

type statusUsecase struct {
	tm TxManager
	sr repository.StatusRepository
}

func NewStatusUsecase(tm TxManager, sr repository.StatusRepository) *statusUsecase {
	return &statusUsecase{tm: tm, sr: sr}
}

func (su *statusUsecase) GetOrCreate(ctx context.Context, status *domain.Status) (*domain.Status, error) {
	return su.sr.GetOrCreate(ctx, status)
}
