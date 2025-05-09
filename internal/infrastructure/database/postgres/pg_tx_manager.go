package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sorrawichYooboon/online-order-management-service/internal/repository"
)

type PgxTxManager struct {
	db *pgxpool.Pool
}

func NewTxManager(db *pgxpool.Pool) repository.PgTxManager {
	return &PgxTxManager{db: db}
}

func (t *PgxTxManager) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	conn, err := t.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
