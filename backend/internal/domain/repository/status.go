package repository

import "backend/internal/domain/entity"

type StatusRepository interface {
	GetOrCreate(status *entity.Status) (*entity.Status, error)
}
