package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
	"github.com/sorrawichYooboon/online-order-management-service/internal/repository"
)

type OrderRepositoryImpl struct {
	db *pgx.Conn
}

func NewOrderRepository(db *pgx.Conn) repository.OrderRepository {
	return &OrderRepositoryImpl{
		db: db,
	}
}

func (r *OrderRepositoryImpl) InsertOrder(ctx context.Context, tx pgx.Tx, order *domain.Order) (int, error) {
	var id int
	err := tx.QueryRow(ctx,
		`INSERT INTO orders (customer_name, total_amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		order.CustomerName, order.TotalAmount, order.Status, order.CreatedAt, order.UpdatedAt,
	).Scan(&id)
	return id, err
}
