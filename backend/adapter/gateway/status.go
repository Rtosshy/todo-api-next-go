package gateway

import (
	"backend/entity"
	"backend/usecase"

	"gorm.io/gorm"
)

type statusRepository struct {
	db *gorm.DB
}

func NewStatusRepository(db *gorm.DB) usecase.IStatusRepository {
	return &statusRepository{db: db}
}

func (sr *statusRepository) GetOrCreateStatus(status *entity.Status) (*entity.Status, error) {
	var getOrCreateStatus entity.Status
	if err := sr.db.FirstOrCreate(&getOrCreateStatus, status).Error; err != nil {
		return nil, err
	}
	return &getOrCreateStatus, nil
}
