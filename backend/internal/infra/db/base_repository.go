package db

import (
	"context"
	"gorm.io/gorm"
)

type baseRepository struct {
	_db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) baseRepository {
	return baseRepository{_db: db}
}

func (b *baseRepository) db(ctx context.Context) *gorm.DB {
	if db, ok := GetTx(ctx); ok {
		return db
	}
	return b._db.WithContext(ctx)
}
