package apperror

func NewInternalError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    500,
	}
}
