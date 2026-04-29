package usecase

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"context"
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
