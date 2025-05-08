package usecase

type OrderUsecaseImpl struct {
}

func NewOrderUsecase() OrderUsecase {
	return &OrderUsecaseImpl{}
}
