package apperror

type UseCaseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *UseCaseError) Error() string {
	return e.Message
}

var (
	ErrDatabase   = UseCaseError{Code: "ERROR_DATABASE", Message: "Database error"}
	ErrUnexpected = UseCaseError{Code: "ERROR_UNEXPECTED", Message: "Unexpected error"}
)
