package repository

import (
	"backend/internal/domain"
	"context"
)

type StatusRepository interface {
	GetOrCreate(ctx context.Context, status *domain.Status) (*domain.Status, error)
}
