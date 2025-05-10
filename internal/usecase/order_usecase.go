package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
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

func (u *OrderUsecaseImpl) CreateOrder(ctx context.Context, orders []domain.Order) error {
	var resultErr error
	errChan := make(chan error, len(orders))

	for _, o := range orders {
		order := o
		u.workerPool.AddTask(workers.Task{
			Execute: func() {
				defer func() {
					if r := recover(); r != nil {
						errChan <- fmt.Errorf("panic: %v", r)
					}
				}()

				localCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

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

					for i := range order.Items {
						order.Items[i].OrderID = orderID
					}

					return u.orderItemRepo.InsertBatchTx(localCtx, tx, order.Items)
				})

				errChan <- err
			},
		})
	}

	for range orders {
		if err := <-errChan; err != nil {
			resultErr = multierror.Append(resultErr, err)
			log.Printf("Error processing order: %v", err)
		}
	}

	return resultErr
}

func (u *OrderUsecaseImpl) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	return u.pgTxManager.WithTx(ctx, func(tx pgx.Tx) error {
		if err := u.orderRepo.UpdateStatusTx(ctx, tx, orderID, status); err != nil {
			return err
		}

		return nil
	})
}
