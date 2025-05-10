package rest

import (
	"github.com/labstack/echo/v4"

	"github.com/sorrawichYooboon/online-order-management-service/config"
	"github.com/sorrawichYooboon/online-order-management-service/internal/infrastructure/delivery/rest/handler"
)

func NewServer(e *echo.Echo, cfg *config.Config, healthHandler handler.HealthHandler, orderHandler handler.OrderHandler) {
	e.GET("/ping", healthHandler.Ping)

	orders := e.Group("/orders")
	orders.GET("", orderHandler.GetOrders)
	orders.GET("/:order_id", orderHandler.GetOrderByID)
	orders.POST("", orderHandler.CreateOrders)
	orders.PUT("/:order_id/status", orderHandler.UpdateOrderStatus)
}
