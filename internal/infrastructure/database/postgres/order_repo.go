package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
	"github.com/sorrawichYooboon/online-order-management-service/internal/repository"
)

type OrderRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) repository.OrderRepository {
	return &OrderRepositoryImpl{
		db: db,
	}
}

func (r *OrderRepositoryImpl) GetPaginated(ctx context.Context, page int, pageSize int, sort string) ([]domain.Order, error) {
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, customer_name, status, total_amount, created_at, updated_at
		FROM orders
		ORDER BY created_at %s
		LIMIT $1 OFFSET $2
	`, sort)

	rows, err := r.db.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.CustomerName, &o.Status, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *OrderRepositoryImpl) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, customer_name, status, total_amount, created_at, updated_at
		FROM orders
		WHERE id = $1
	`, id)

	var o domain.Order
	if err := row.Scan(&o.ID, &o.CustomerName, &o.Status, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt); err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, order_id, product_name, quantity, price
		FROM order_items
		WHERE order_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductName, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}

	return &o, nil
}

func (r *OrderRepositoryImpl) InsertTx(ctx context.Context, tx pgx.Tx, order *domain.Order) (int, error) {
	var id int
	err := tx.QueryRow(ctx,
		`INSERT INTO orders (customer_name, total_amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		order.CustomerName, order.TotalAmount, order.Status, order.CreatedAt, order.UpdatedAt,
	).Scan(&id)
	return id, err
}
