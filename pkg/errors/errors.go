package errors

import "net/http"

type AppError struct {
	HTTPStatus int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}


func New(status int, message string) *AppError {
	return &AppError{HTTPStatus: status, Message: message}
}

func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message)
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, message)
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, message)
}

func Conflict(message string) *AppError {
	return New(http.StatusConflict, message)
}

func UnprocessableEntity(message string) *AppError {
	return New(http.StatusUnprocessableEntity, message)
}

func InternalServerError(message string) *AppError {
	return New(http.StatusInternalServerError, message)
}


var (
	ErrTransactionNotFound     = NotFound("transaction not found")
	ErrAccountNotFound         = NotFound("account not found")
	ErrUserNotFound            = NotFound("user not found")
	ErrInsufficientBalance     = BadRequest("insufficient balance")
	ErrAccountInactive         = BadRequest("account is inactive")
	ErrSameAccount             = BadRequest("source and destination account cannot be the same")
	ErrInvalidTransactionType  = BadRequest("invalid transaction type")
	ErrTransactionAlreadyFinal = Conflict("transaction is already in a final state")
)