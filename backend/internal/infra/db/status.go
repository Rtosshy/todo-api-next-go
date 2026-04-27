package db

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"

	"gorm.io/gorm"
)

type statusRepository struct {
	db *gorm.DB
}

func NewStatusRepository(db *gorm.DB) repository.StatusRepository {
	return &statusRepository{db: db}
}

func (sr *statusRepository) GetOrCreate(status *domain.Status) (*domain.Status, error) {
	var getOrCreateStatus domain.Status
	if err := sr.db.FirstOrCreate(&getOrCreateStatus, status).Error; err != nil {
		return nil, err
	}
	return &getOrCreateStatus, nil
}
