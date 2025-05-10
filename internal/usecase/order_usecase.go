package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
	"github.com/sorrawichYooboon/online-order-management-service/internal/repository"
	"github.com/sorrawichYooboon/online-order-management-service/logger"
	"github.com/sorrawichYooboon/online-order-management-service/pkg/apperror"
	"github.com/sorrawichYooboon/online-order-management-service/pkg/retry"
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
	orders, err := u.orderRepo.GetPaginated(ctx, page, pageSize, sort)
	if err != nil {
		logger.LogError(ORDER_USECASE_GET_ORDERS, err)
		return nil, &apperror.ErrDatabase
	}

	return orders, nil
}

func (u *OrderUsecaseImpl) GetOrderByID(ctx context.Context, id int64) (*domain.Order, error) {
	order, err := u.orderRepo.GetByID(ctx, id)
	if err != nil {
		logger.LogError(ORDER_USECASE_GET_ORDER_BY_ID, err)
		return nil, &apperror.ErrDatabase
	}

	return order, nil
}

func (u *OrderUsecaseImpl) CreateOrders(ctx context.Context, orders []domain.Order) ([]CreateOrdersResponse, error) {
	resultChan := make(chan CreateOrdersResponse, len(orders))

	for index, o := range orders {
		order := o
		i := index

		u.workerPool.AddTask(workers.Task{
			Execute: func() {
				defer func() {
					if r := recover(); r != nil {
						resultChan <- CreateOrdersResponse{
							Index: i,
							Error: fmt.Errorf("panic: %v", r),
						}
					}
				}()

				localCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				var result CreateOrdersResponse
				result.Index = i

				err := retry.RetryIf(
					3,
					100*time.Millisecond,
					2*time.Second,
					retry.IsTransientError,
					func() error {
						return u.pgTxManager.WithTx(localCtx, func(tx pgx.Tx) error {
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
					})

				if err != nil {
					logger.LogError(ORDER_USECASE_CREATE_ORDER, err)
					result.Error = &apperror.ErrDatabase
				}

				resultChan <- result
			},
		})
	}

	results := make([]CreateOrdersResponse, len(orders))
	for range orders {
		r := <-resultChan
		results[r.Index] = r
	}

	return results, nil
}

func (u *OrderUsecaseImpl) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	return u.pgTxManager.WithTx(ctx, func(tx pgx.Tx) error {
		if err := u.orderRepo.UpdateStatusTx(ctx, tx, orderID, status); err != nil {
			logger.LogError(ORDER_USECASE_UPDATE_ORDER, err)
			return &apperror.ErrDatabase
		}

		return nil
	})
}
