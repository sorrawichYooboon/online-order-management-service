package dto

import "github.com/sorrawichYooboon/online-order-management-service/internal/domain"

type GetOrdersRequestDTO struct {
	Page     int    `query:"page" validate:"min=1"`
	PageSize int    `query:"page_size" validate:"min=1,max=100"`
	Sort     string `query:"sort"`
}

type GetOrdersResponseDTO struct {
	Summary GetOrdersSummaryDTO `json:"summary"`
	Orders  []domain.Order      `json:"orders"`
}

type GetOrdersSummaryDTO struct {
	Page              int `json:"page"`
	PageSize          int `json:"page_size"`
	TotalOrdersOnPage int `json:"total_orders_on_page"`
	TotalItems        int `json:"total_items"`
	TotalPages        int `json:"total_pages"`
}

type CreateOrderRequestDTO struct {
	CustomerName string               `json:"customer_name" validate:"required" example:"Alice Smith"`
	Items        []CreateOrderItemDTO `json:"items" validate:"required,dive"`
}

type CreateOrderItemDTO struct {
	ProductName string  `json:"product_name" validate:"required" example:"Wireless Mouse"`
	Quantity    int     `json:"quantity" validate:"required,gt=0" example:"2"`
	Price       float64 `json:"price" validate:"required,gte=0" example:"19.99"`
}

type CreateOrdersResponseDTO struct {
	Summary OrderInsertSummaryDTO  `json:"summary"`
	Results []OrderInsertResultDTO `json:"results"`
}

type OrderInsertSummaryDTO struct {
	Total   int `json:"total"`
	Success int `json:"success"`
	Failed  int `json:"failed"`
}

type OrderInsertResultDTO struct {
	Index   int    `json:"index"`
	OrderID int    `json:"order_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

type UpdateOrderStatusRequestDTO struct {
	Status string `json:"status" validate:"required,oneof=PENDING PAID SHIPPED CANCELED" example:"PAID"`
}
