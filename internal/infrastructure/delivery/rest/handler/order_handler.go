package handler

import (
	"fmt"
	"net/http"
	"strconv"

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

func (oh *OrderHandlerImpl) GetOrders(c echo.Context) error {
	var req dto.GetOrdersRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid query params"})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if req.Sort != "asc" && req.Sort != "desc" {
		req.Sort = "desc"
	}

	orders, err := oh.orderUsecase.GetOrders(c.Request().Context(), req.Page, req.PageSize, req.Sort)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders"})
	}
	return c.JSON(http.StatusOK, orders)
}

func (oh *OrderHandlerImpl) GetOrderByID(c echo.Context) error {
	idParam := c.Param("order_id")
	orderID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order ID"})
	}

	order, err := oh.orderUsecase.GetOrderByID(c.Request().Context(), orderID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found"})
	}

	return c.JSON(http.StatusOK, order)
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

func (oh *OrderHandlerImpl) UpdateOrderStatus(c echo.Context) error {
	orderIDParam := c.Param("order_id")
	orderID, err := strconv.ParseInt(orderIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order ID"})
	}

	var req dto.UpdateOrderStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := oh.orderUsecase.UpdateOrderStatus(c.Request().Context(), orderID, req.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update order status"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Order status updated successfully"})
}
