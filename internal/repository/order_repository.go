package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
)

type OrderRepository interface {
	GetPaginated(ctx context.Context, page int, pageSize int, sort string) ([]domain.Order, error)
	GetByID(ctx context.Context, id int64) (*domain.Order, error)
	InsertTx(ctx context.Context, tx pgx.Tx, order *domain.Order) (int, error)
}
