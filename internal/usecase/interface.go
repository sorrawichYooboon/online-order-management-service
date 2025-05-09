package usecase

import (
	"context"

	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
)

type OrderUsecase interface {
	GetOrders(ctx context.Context, page, pageSize int, sort string) ([]domain.Order, error)
	GetOrderByID(ctx context.Context, id int64) (*domain.Order, error)
	CreateOrder(ctx context.Context, order []domain.Order) error
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) error
}
