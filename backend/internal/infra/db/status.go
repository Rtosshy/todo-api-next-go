package db

import (
	"backend/internal/domain/entity"
	"backend/internal/domain/repository"

	"gorm.io/gorm"
)

type statusRepository struct {
	db *gorm.DB
}

func NewStatusRepository(db *gorm.DB) repository.StatusRepository {
	return &statusRepository{db: db}
}

func (sr *statusRepository) GetOrCreate(status *entity.Status) (*entity.Status, error) {
	var getOrCreateStatus entity.Status
	if err := sr.db.FirstOrCreate(&getOrCreateStatus, status).Error; err != nil {
		return nil, err
	}
	return &getOrCreateStatus, nil
}
