package db

import (
	"context"
	"gorm.io/gorm"
)

type txKeyType struct{}

var txKey txKeyType

func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTx(ctx context.Context) (*gorm.DB, bool) {
	db, ok := ctx.Value(txKey).(*gorm.DB)
	return db, ok
}
