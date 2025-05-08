package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PgTxManager interface {
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}
