package httperror

type HTTPError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

var (
	ErrInternalServer = HTTPError{Code: "INTERNAL_ERROR", Message: "Internal server error"}
)
