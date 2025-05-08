package dto

type CreateOrderRequest struct {
	CustomerName string               `json:"customer_name" validate:"required"`
	Status       string               `json:"status" validate:"required"`
	Items        []CreateOrderItemDTO `json:"items" validate:"required,dive"`
}

type CreateOrderItemDTO struct {
	ProductName string  `json:"product_name" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,gt=0"`
	Price       float64 `json:"price" validate:"required,gte=0"`
}
