package db

import (
	"backend/internal/usecase"
	"context"
	"gorm.io/gorm"
)

type txManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) usecase.TxManager {
	return &txManager{db: db}
}

func (tm *txManager) RunInTx(ctx context.Context, fn usecase.TxFunc) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctxWithTx := WithTx(ctx, tx)
		if err := fn(ctxWithTx); err != nil {
			return err
		}
		return nil
	})

}
