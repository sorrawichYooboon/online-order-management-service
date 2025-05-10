package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
	"github.com/sorrawichYooboon/online-order-management-service/internal/dto"
	"github.com/sorrawichYooboon/online-order-management-service/internal/usecase"
	"github.com/sorrawichYooboon/online-order-management-service/logger"
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
	var req dto.GetOrdersRequestDTO
	if err := c.Bind(&req); err != nil {
		logger.LogError(ORDER_HANDLER_GET_ORDERS, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid query params"})
	}
	if err := c.Validate(req); err != nil {
		logger.LogError(ORDER_HANDLER_GET_ORDERS, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if req.Sort != "asc" && req.Sort != "desc" {
		req.Sort = "desc"
	}

	orders, err := oh.orderUsecase.GetOrders(c.Request().Context(), req.Page, req.PageSize, req.Sort)
	if err != nil {
		logger.LogError(ORDER_HANDLER_GET_ORDERS, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders"})
	}

	return c.JSON(http.StatusOK, orders)
}

func (oh *OrderHandlerImpl) GetOrderByID(c echo.Context) error {
	idParam := c.Param("order_id")
	orderID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		logger.LogError(ORDER_HANDLER_GET_ORDER_BY_ID, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order ID"})
	}

	order, err := oh.orderUsecase.GetOrderByID(c.Request().Context(), orderID)
	if err != nil {
		logger.LogError(ORDER_HANDLER_GET_ORDER_BY_ID, err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found"})
	}

	return c.JSON(http.StatusOK, order)
}

func (oh *OrderHandlerImpl) CreateOrders(c echo.Context) error {
	var createOrderRequest []dto.CreateOrderRequestDTO
	if err := c.Bind(&createOrderRequest); err != nil {
		logger.LogError(ORDER_HANDLER_CREATE_ORDER, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	results := make([]dto.OrderInsertResultDTO, len(createOrderRequest))
	validOrders := make([]domain.Order, 0, len(createOrderRequest))
	validIndexes := make([]int, 0, len(createOrderRequest))

	for i, orderRequest := range createOrderRequest {
		if err := c.Validate(orderRequest); err != nil {
			logger.LogError(ORDER_HANDLER_CREATE_ORDER, err)
			results[i] = dto.OrderInsertResultDTO{
				Index: i,
				Error: fmt.Sprintf("validation failed: %v", err),
			}
			continue
		}

		order := domain.Order{
			CustomerName: orderRequest.CustomerName,
			Status:       orderRequest.Status,
			Items:        make([]domain.OrderItem, len(orderRequest.Items)),
		}
		for j, item := range orderRequest.Items {
			order.Items[j] = domain.OrderItem{
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				Price:       item.Price,
			}
		}

		validOrders = append(validOrders, order)
		validIndexes = append(validIndexes, i)
	}

	if len(validOrders) > 0 {
		usecaseResults, _ := oh.orderUsecase.CreateOrders(c.Request().Context(), validOrders)

		for j, r := range usecaseResults {
			originalIndex := validIndexes[j]
			results[originalIndex] = dto.OrderInsertResultDTO{
				Index:   originalIndex,
				OrderID: r.OrderID,
			}
			if r.Error != "" {
				results[originalIndex].Error = r.Error
			}
		}
	}
	summary := dto.OrderInsertSummaryDTO{}
	for _, r := range results {
		summary.Total++
		if r.Error != "" {
			summary.Failed++
		} else {
			summary.Success++
		}
	}

	status := http.StatusCreated
	if summary.Success == 0 {
		status = http.StatusBadRequest

	} else if summary.Failed > 0 {
		status = http.StatusOK

	}

	return c.JSON(status, dto.CreateOrderResponseDTO{
		Summary: summary,
		Results: results,
	})
}

func (oh *OrderHandlerImpl) UpdateOrderStatus(c echo.Context) error {
	orderIDParam := c.Param("order_id")
	orderID, err := strconv.ParseInt(orderIDParam, 10, 64)
	if err != nil {
		logger.LogError(ORDER_HANDLER_UPDATE_ORDER_STATUS, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order ID"})
	}

	var req dto.UpdateOrderStatusRequestDTO
	if err := c.Bind(&req); err != nil {
		logger.LogError(ORDER_HANDLER_UPDATE_ORDER_STATUS, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	if err := c.Validate(req); err != nil {
		logger.LogError(ORDER_HANDLER_UPDATE_ORDER_STATUS, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := oh.orderUsecase.UpdateOrderStatus(c.Request().Context(), orderID, req.Status); err != nil {
		logger.LogError(ORDER_HANDLER_UPDATE_ORDER_STATUS, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update order status"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Order status updated successfully"})
}
