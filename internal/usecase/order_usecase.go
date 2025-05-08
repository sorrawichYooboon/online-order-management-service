package usecase

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
	"github.com/sorrawichYooboon/online-order-management-service/internal/repository"
)

type OrderUsecaseImpl struct {
	pgTxManager   repository.PgTxManager
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
}

func NewOrderUsecase(
	pgTxManager repository.PgTxManager,
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
) OrderUsecase {
	return &OrderUsecaseImpl{
		pgTxManager:   pgTxManager,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
	}
}

func (u *OrderUsecaseImpl) CreateOrder(ctx context.Context, orders []domain.Order) error {
	for _, order := range orders {
		err := u.pgTxManager.WithTx(ctx, func(tx pgx.Tx) error {
			order.CreatedAt = time.Now()
			order.UpdatedAt = time.Now()

			var total float64
			for _, item := range order.Items {
				total += float64(item.Quantity) * item.Price
			}
			order.TotalAmount = total

			orderID, err := u.orderRepo.InsertOrder(ctx, tx, &order)
			if err != nil {
				return err
			}

			for i := range order.Items {
				order.Items[i].OrderID = orderID
			}

			return u.orderItemRepo.InsertOrderItems(ctx, tx, order.Items)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
