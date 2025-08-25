package apperror

func NewValidationError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    400,
	}
}
