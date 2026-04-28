package repository

import (
	"backend/internal/domain"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Get(ctx context.Context, userID domain.UserID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, userID domain.UserID) error
}
