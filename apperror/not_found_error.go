package apperror

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    404,
	}
}
