package apperror

type AppError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *AppError) Error() string {
	return e.Message
}
