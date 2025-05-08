package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthHandlerImpl struct {
}

func NewHealthHandler() HealthHandler {
	return &HealthHandlerImpl{}
}

func (h *HealthHandlerImpl) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Pong~!"})
}
