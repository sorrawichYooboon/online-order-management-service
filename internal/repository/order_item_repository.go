package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
)

type OrderItemRepository interface {
	InsertBatchTx(ctx context.Context, tx pgx.Tx, items []domain.OrderItem) error
}
