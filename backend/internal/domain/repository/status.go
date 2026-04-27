package repository

import "backend/internal/domain"

type StatusRepository interface {
	GetOrCreate(status *domain.Status) (*domain.Status, error)
}
