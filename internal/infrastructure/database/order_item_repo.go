package database

import "github.com/sorrawichYooboon/online-order-management-service/internal/repository"

type OrderItemRepositoryImpl struct {
}

func NewOrderItemRepository() repository.OrderItemRepository {
	return &OrderItemRepositoryImpl{}
}
