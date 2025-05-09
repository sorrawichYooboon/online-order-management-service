package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
	"github.com/sorrawichYooboon/online-order-management-service/internal/repository"
)

type OrderItemRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewOrderItemRepository(db *pgxpool.Pool) repository.OrderItemRepository {
	return &OrderItemRepositoryImpl{
		db: db,
	}
}

func (r *OrderItemRepositoryImpl) InsertOrderItems(ctx context.Context, tx pgx.Tx, items []domain.OrderItem) error {
	batch := &pgx.Batch{}

	for _, item := range items {
		batch.Queue(`
			INSERT INTO order_items (order_id, product_name, quantity, price)
			VALUES ($1, $2, $3, $4)
		`, item.OrderID, item.ProductName, item.Quantity, item.Price)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for range items {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}
