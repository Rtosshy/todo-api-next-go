package gateway

import (
	"backend/internal/domain"
	"backend/internal/domain/repo"

	"gorm.io/gorm"
)

type statusRepositoryImpl struct {
	db *gorm.DB
}

func NewStatusRepository(db *gorm.DB) repo.StatusRepo {
	return &statusRepositoryImpl{db: db}
}

func (sr *statusRepositoryImpl) GetOrCreate(status *domain.Status) (*domain.Status, error) {
	var getOrCreateStatus domain.Status
	if err := sr.db.FirstOrCreate(&getOrCreateStatus, status).Error; err != nil {
		return nil, err
	}
	return &getOrCreateStatus, nil
}
