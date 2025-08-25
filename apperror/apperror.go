package apperror

// AppError defines a custom error
type AppError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}
