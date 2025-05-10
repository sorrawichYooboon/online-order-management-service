package response

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var (
	SuccessPing = Response{Code: "SUCCESS_PING", Message: "Pong~!"}

	SuccessOrderGetOrders    = Response{Code: "SUCCESS_GET_ORDERS", Message: "Get orders successfully"}
	SuccessOrderGetOrderByID = Response{Code: "SUCCESS_GET_ORDER_BY_ID", Message: "Get order by ID successfully"}
	SuccessOrderCreateOrders = Response{Code: "SUCCESS_CREATE_ORDERS", Message: "Create orders successfully"}
	SuccessOrderUpdateStatus = Response{Code: "SUCCESS_UPDATE_ORDER_STATUS", Message: "Update order status successfully"}
)
