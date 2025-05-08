package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func PingHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Pong!",
	})
}
