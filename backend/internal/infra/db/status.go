package db

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"backend/internal/infra/db/dao"

	"context"
)

type statusRepository struct {
	baseRepository
}

func NewStatusRepository(b baseRepository) repository.StatusRepository {
	return &statusRepository{b}
}

func (sr *statusRepository) GetOrCreate(ctx context.Context, status *domain.Status) (*domain.Status, error) {
	want := dao.Status{
		ID:   dao.StatusID(status.ID),
		Name: dao.StatusName(status.Name.String()),
	}
	var got dao.Status
	if err := sr.db(ctx).FirstOrCreate(&got, want).Error; err != nil {
		return nil, err
	}
	statusName, err := domain.NewStatusName(string(got.Name))
	if err != nil {
		return nil, err
	}
	return &domain.Status{ID: domain.StatusID(got.ID), Name: statusName}, nil
}
