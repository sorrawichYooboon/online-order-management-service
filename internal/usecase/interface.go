package usecase

import (
	"context"

	"github.com/sorrawichYooboon/online-order-management-service/internal/domain"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, order []domain.Order) error
}
