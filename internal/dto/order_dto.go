package dto

type GetOrdersRequestDTO struct {
	Page     int    `query:"page" validate:"min=1"`
	PageSize int    `query:"page_size" validate:"min=1,max=100"`
	Sort     string `query:"sort"`
}

type CreateOrderRequestDTO struct {
	CustomerName string               `json:"customer_name" validate:"required"`
	Status       string               `json:"status" validate:"required"`
	Items        []CreateOrderItemDTO `json:"items" validate:"required,dive"`
}

type CreateOrderItemDTO struct {
	ProductName string  `json:"product_name" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,gt=0"`
	Price       float64 `json:"price" validate:"required,gte=0"`
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
	Status string `json:"status" validate:"required,oneof=PENDING PAID SHIPPED CANCELED"`
}
