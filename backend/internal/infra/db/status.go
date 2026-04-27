package db

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"backend/internal/infra/db/dao"

	"gorm.io/gorm"
)

type statusRepository struct {
	db *gorm.DB
}

func NewStatusRepository(db *gorm.DB) repository.StatusRepository {
	return &statusRepository{db: db}
}

func (sr *statusRepository) GetOrCreate(status *domain.Status) (*domain.Status, error) {
	want := dao.Status{
		ID:   dao.StatusID(status.ID),
		Name: dao.StatusName(status.Name.String()),
	}
	var got dao.Status
	if err := sr.db.FirstOrCreate(&got, want).Error; err != nil {
		return nil, err
	}
	statusName, err := domain.NewStatusName(string(got.Name))
	if err != nil {
		return nil, err
	}
	return &domain.Status{ID: domain.StatusID(got.ID), Name: statusName}, nil
}
