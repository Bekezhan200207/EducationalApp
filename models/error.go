package models

type ApiError struct {
	Error string
}

func NewApiError(message string) ApiError {
	return ApiError{message}
}
