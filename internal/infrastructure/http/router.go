package http

import (
	"github.com/labstack/echo/v4"
	"github.com/sorrawichYooboon/online-order-management-service/internal/infrastructure/http/handler"
)

func NewEchoServer() *echo.Echo {
	e := echo.New()

	// Register routes
	e.GET("/ping", handler.PingHandler)

	return e
}
