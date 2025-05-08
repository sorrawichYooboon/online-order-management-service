package handler

type OrderHandlerImpl struct {
}

func NewOrderHandler() OrderHandler {
	return &OrderHandlerImpl{}
}
