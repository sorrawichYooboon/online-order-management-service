package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
	"github.com/sorrawichYooboon/online-order-management-service/internal/dto"
	"github.com/sorrawichYooboon/online-order-management-service/internal/usecase"
)

type OrderHandlerImpl struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase usecase.OrderUsecase) OrderHandler {
	return &OrderHandlerImpl{
		orderUsecase: orderUsecase,
	}
}

func (oh *OrderHandlerImpl) CreateOrder(c echo.Context) error {
	var createOrderRequest []dto.CreateOrderRequest
	if err := c.Bind(&createOrderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	orders := make([]domain.Order, len(createOrderRequest))
	for i, orderRequest := range createOrderRequest {

		if err := c.Validate(createOrderRequest[i]); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("validation failed on order %d: %v", i+1, err),
			})
		}

		orders[i] = domain.Order{
			CustomerName: orderRequest.CustomerName,
			Status:       orderRequest.Status,
			Items:        make([]domain.OrderItem, len(orderRequest.Items)),
		}
		for j, item := range orderRequest.Items {
			orders[i].Items[j] = domain.OrderItem{
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				Price:       item.Price,
			}
		}
	}

	if err := oh.orderUsecase.CreateOrder(c.Request().Context(), orders); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create order"})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "Order created successfully"})
}
