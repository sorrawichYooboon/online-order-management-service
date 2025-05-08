package apperror

type UseCaseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *UseCaseError) Error() string {
	return e.Message
}

var ()
