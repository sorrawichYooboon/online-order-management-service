package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sorrawichYooboon/online-order-management-service/internal/repository"
)

type PgxTxManager struct {
	DB *pgx.Conn
}

func NewTxManager(db *pgx.Conn) repository.PgTxManager {
	return &PgxTxManager{DB: db}
}

func (t *PgxTxManager) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := t.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
