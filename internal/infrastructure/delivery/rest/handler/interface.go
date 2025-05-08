package handler

import "github.com/labstack/echo/v4"

type HealthHandler interface {
	Ping(echo.Context) error
}

type OrderHandler interface {
	CreateOrder(echo.Context) error
}
