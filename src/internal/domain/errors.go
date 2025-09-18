package domain

import "fmt"

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
	Query   string `json:"query,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewValidationError(message, query string) *AppError {
	return &AppError{
		Code:    400,
		Message: message,
		Query:   query,
	}
}

func NewNotFoundError(query string) *AppError {
	return &AppError{
		Code:    404,
		Message: fmt.Sprintf("Country not found: %s", query),
		Query:   query,
	}
}

func NewInternalError(message string) *AppError {
	return &AppError{
		Code:    500,
		Message: message,
	}
}
