package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sorrawichYooboon/online-order-management-service/pkg/response"
)

type HealthHandlerImpl struct {
}

func NewHealthHandler() HealthHandler {
	return &HealthHandlerImpl{}
}

func (h *HealthHandlerImpl) Ping(c echo.Context) error {
	return response.Success(c, http.StatusOK, response.SuccessPing, nil)
}
