package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
	"github.com/sorrawichYooboon/online-order-management-service/internal/repository"
	"github.com/sorrawichYooboon/online-order-management-service/pkg/workers"
)

type OrderUsecaseImpl struct {
	pgTxManager   repository.PgTxManager
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
	workerPool    *workers.WorkerPool
}

func NewOrderUsecase(
	pgTxManager repository.PgTxManager,
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	workerPool *workers.WorkerPool,
) OrderUsecase {
	return &OrderUsecaseImpl{
		pgTxManager:   pgTxManager,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		workerPool:    workerPool,
	}
}

func (u *OrderUsecaseImpl) GetOrders(ctx context.Context, page int, pageSize int, sort string) ([]domain.Order, error) {
	return u.orderRepo.GetPaginated(ctx, page, pageSize, sort)
}

func (u *OrderUsecaseImpl) GetOrderByID(ctx context.Context, id int64) (*domain.Order, error) {
	return u.orderRepo.GetByID(ctx, id)
}

func (u *OrderUsecaseImpl) CreateOrder(ctx context.Context, orders []domain.Order) ([]CreateOrderResponse, error) {
	resultChan := make(chan CreateOrderResponse, len(orders))

	for index, o := range orders {
		order := o
		i := index

		u.workerPool.AddTask(workers.Task{
			Execute: func() {
				defer func() {
					if r := recover(); r != nil {
						resultChan <- CreateOrderResponse{
							Index: i,
							Error: fmt.Sprintf("panic: %v", r),
						}
					}
				}()

				localCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				var result CreateOrderResponse
				result.Index = i

				err := u.pgTxManager.WithTx(localCtx, func(tx pgx.Tx) error {
					var total float64
					for _, item := range order.Items {
						total += float64(item.Quantity) * item.Price
					}
					order.TotalAmount = total

					orderID, err := u.orderRepo.InsertTx(localCtx, tx, &order)
					if err != nil {
						return err
					}

					for j := range order.Items {
						order.Items[j].OrderID = orderID
					}

					if err := u.orderItemRepo.InsertBatchTx(localCtx, tx, order.Items); err != nil {
						return err
					}

					result.OrderID = orderID
					return nil
				})

				if err != nil {
					result.Error = err.Error()
				}

				resultChan <- result
			},
		})
	}

	results := make([]CreateOrderResponse, len(orders))
	for range orders {
		r := <-resultChan
		results[r.Index] = r
	}

	return results, nil
}

func (u *OrderUsecaseImpl) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	return u.pgTxManager.WithTx(ctx, func(tx pgx.Tx) error {
		if err := u.orderRepo.UpdateStatusTx(ctx, tx, orderID, status); err != nil {
			return err
		}

		return nil
	})
}
