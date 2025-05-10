package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderRepository interface {
	GetPaginated(ctx context.Context, page int, pageSize int, sort string) ([]domain.Order, int, error)
	GetByID(ctx context.Context, id int64) (*domain.Order, error)
	InsertTx(ctx context.Context, tx pgx.Tx, order *domain.Order) (int, error)
	UpdateStatusTx(ctx context.Context, tx pgx.Tx, orderID int64, status string) error
}
