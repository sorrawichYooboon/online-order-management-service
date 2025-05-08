package rest

import (
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"github.com/sorrawichYooboon/online-order-management-service/config"
	"github.com/sorrawichYooboon/online-order-management-service/internal/infrastructure/delivery/rest/handler"
)

func NewServer(e *echo.Echo, cfg *config.Config, pq *pgx.Conn, healthHandler handler.HealthHandler, orderHandler handler.OrderHandler) {
	e.GET("/ping", healthHandler.Ping)

	orders := e.Group("/orders")
	orders.POST("", orderHandler.CreateOrder)
}
