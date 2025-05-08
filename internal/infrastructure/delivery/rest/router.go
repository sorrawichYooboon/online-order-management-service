package rest

import (
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/sorrawichYooboon/online-order-management-service/config"
	"github.com/sorrawichYooboon/online-order-management-service/internal/infrastructure/delivery/rest/handler"
)

func NewServer(*config.Config, *pgx.Conn) *echo.Echo {
	e := echo.New()

	e.GET("/ping", handler.PingHandler)

	return e
}
