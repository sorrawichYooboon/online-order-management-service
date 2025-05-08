package domain

import "time"

type Order struct {
	ID           int
	CustomerName string
	TotalAmount  float64
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Items        []OrderItem
}
