package util

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type ErrorCode int

const (
	ErrInternalServerErrorCode ErrorCode = iota
	ErrNotFoundCode
	ErrConflictCode
	ErrUnauthorizedCode
	ErrBadRequestCode
)

type CustomError struct {
	Code    ErrorCode
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewCustomError(code int, message string) *CustomError {
	return &CustomError{
		Code:    ErrorCode(code),
		Message: message,
	}
}

var (
	ErrInternalServerError = &CustomError{Code: ErrInternalServerErrorCode}
	ErrNotFound            = &CustomError{Code: ErrNotFoundCode}
	ErrConflict            = &CustomError{Code: ErrConflictCode}
	ErrUnauthorized        = &CustomError{Code: ErrUnauthorizedCode}
	ErrBadRequest          = &CustomError{Code: ErrUnauthorizedCode}
)

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)

	switch e := err.(type) {
	case *CustomError:
		switch e.Code {
		case ErrInternalServerErrorCode:
			return http.StatusInternalServerError
		case ErrNotFoundCode:
			return http.StatusNotFound
		case ErrConflictCode:
			return http.StatusConflict
		case ErrUnauthorizedCode:
			return http.StatusUnauthorized
		case ErrBadRequestCode:
			return http.StatusBadRequest
		default:
			return http.StatusInternalServerError
		}
	default:
		logrus.Warn("Received an unknown error type")
		return http.StatusInternalServerError
	}
}
