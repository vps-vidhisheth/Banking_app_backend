package apperror

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    401,
	}
}
