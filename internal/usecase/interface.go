package usecase

import (
	"context"

	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
)

type OrderUsecase interface {
	GetOrders(ctx context.Context, page int, pageSize int, sort string) ([]domain.Order, int, error)
	GetOrderByID(ctx context.Context, id int64) (*domain.Order, error)
	CreateOrders(ctx context.Context, orders []domain.Order) ([]CreateOrdersResponse, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) error
}
