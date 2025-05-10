package apperror

type UseCaseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *UseCaseError) Error() string {
	return e.Message
}

var (
	ErrOrderNotFound = UseCaseError{Code: "ERROR_ORDER_NOT_FOUND", Message: "Order not found"}
	ErrDatabase      = UseCaseError{Code: "ERROR_DATABASE", Message: "Database error"}
	ErrUnexpected    = UseCaseError{Code: "ERROR_UNEXPECTED", Message: "Unexpected error"}
)
