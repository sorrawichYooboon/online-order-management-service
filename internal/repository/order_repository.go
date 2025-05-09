package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
)

type OrderRepository interface {
	InsertTx(ctx context.Context, tx pgx.Tx, order *domain.Order) (int, error)
}
