package usecase

import (
	"context"
)

type TxFunc func(ctx context.Context) error

type TxManager interface {
	RunInTx(ctx context.Context, fn TxFunc) error
}
