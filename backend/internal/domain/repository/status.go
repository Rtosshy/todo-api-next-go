package repository

import (
	"context"

	"backend/internal/domain"
)

type StatusRepository interface {
	GetOrCreate(ctx context.Context, status *domain.Status) (*domain.Status, error)
}
