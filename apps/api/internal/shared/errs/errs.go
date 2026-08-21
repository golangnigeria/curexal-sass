package errs

import (
	"net/http"
)

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequestError(msg string, _ ...interface{}) *AppError {
	return &AppError{
		Code:       "BAD_REQUEST",
		Message:    msg,
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewInternalServerError() *AppError {
	return &AppError{
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "An unexpected error occurred",
		HTTPStatus: http.StatusInternalServerError,
	}
}

func NewNotFoundError(msg string, _ ...interface{}) *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    msg,
		HTTPStatus: http.StatusNotFound,
	}
}

func NewUnauthorizedError(msg string, _ ...interface{}) *AppError {
	return &AppError{
		Code:       "UNAUTHORIZED",
		Message:    msg,
		HTTPStatus: http.StatusUnauthorized,
	}
}

func NewForbiddenError(msg string, _ ...interface{}) *AppError {
	return &AppError{
		Code:       "FORBIDDEN",
		Message:    msg,
		HTTPStatus: http.StatusForbidden,
	}
}

func NewPaymentRequiredError(msg string, _ ...interface{}) *AppError {
	return &AppError{
		Code:       "PAYMENT_REQUIRED",
		Message:    msg,
		HTTPStatus: http.StatusPaymentRequired,
	}
}

var ErrOrganizationNotFound = NewNotFoundError("Organization not found")
