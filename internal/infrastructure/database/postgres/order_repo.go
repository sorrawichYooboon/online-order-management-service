package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (r *OrderRepositoryImpl) GetPaginated(ctx context.Context, page int, pageSize int, sort string) ([]domain.Order, int, error) {
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	offset := (page - 1) * pageSize

	var totalItems int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders`).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, customer_name, status, total_amount, created_at, updated_at
		FROM orders
		ORDER BY id %s
		LIMIT $1 OFFSET $2
	`, sort)

	rows, err := r.db.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	orderMap := make(map[int]*domain.Order)
	var orderIDs []int

	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.CustomerName, &o.Status, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, err
		}
		orderMap[o.ID] = &o
		orderIDs = append(orderIDs, o.ID)
	}

	if len(orderIDs) == 0 {
		return []domain.Order{}, totalItems, nil
	}

	itemRows, err := r.db.Query(ctx, `
		SELECT id, order_id, product_name, quantity, price
		FROM order_items
		WHERE order_id = ANY($1)
	`, orderIDs)
	if err != nil {
		return nil, 0, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var item domain.OrderItem
		if err := itemRows.Scan(&item.ID, &item.OrderID, &item.ProductName, &item.Quantity, &item.Price); err != nil {
			return nil, 0, err
		}
		if order, ok := orderMap[item.OrderID]; ok {
			order.Items = append(order.Items, item)
		}
	}

	var orders []domain.Order
	for _, id := range orderIDs {
		orders = append(orders, *orderMap[id])
	}

	return orders, totalItems, nil
}

func (r *OrderRepositoryImpl) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, customer_name, status, total_amount, created_at, updated_at
		FROM orders
		WHERE id = $1
	`, id)

	var o domain.Order
	if err := row.Scan(&o.ID, &o.CustomerName, &o.Status, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrOrderNotFound
		}
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
	var createdAt, updatedAt time.Time

	err := tx.QueryRow(ctx,
		`INSERT INTO orders (customer_name, total_amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at`,
		order.CustomerName, order.TotalAmount, order.Status,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return 0, err
	}

	order.ID = id
	order.CreatedAt = createdAt
	order.UpdatedAt = updatedAt

	return id, nil
}

func (r *OrderRepositoryImpl) UpdateStatusTx(ctx context.Context, tx pgx.Tx, orderID int64, status string) error {
	cmd, err := tx.Exec(ctx, `
		UPDATE orders SET status = $1, updated_at = NOW()
		WHERE id = $2
	`, status, orderID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return repository.ErrOrderNotFound
	}

	return nil
}
