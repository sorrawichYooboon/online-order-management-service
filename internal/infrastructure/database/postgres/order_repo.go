package postgres

import "github.com/sorrawichYooboon/online-order-management-service/internal/repository"

type OrderRepositoryImpl struct {
}

func NewOrderRepository() repository.OrderRepository {
	return &OrderRepositoryImpl{}
}
