package handler

import "github.com/labstack/echo/v4"

type HealthHandler interface {
	Ping(echo.Context) error
}

type OrderHandler interface {
	GetOrders(echo.Context) error
	GetOrderByID(echo.Context) error
	CreateOrders(echo.Context) error
	UpdateOrderStatus(echo.Context) error
}
