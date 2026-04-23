package repo

import "backend/internal/domain"

type StatusRepo interface {
	GetOrCreate(status *domain.Status) (*domain.Status, error)
}
