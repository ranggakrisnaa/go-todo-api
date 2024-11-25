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
	ErrInternalServerError = &CustomError{Code: ErrInternalServerErrorCode, Message: "Internal Server Error"}
	ErrNotFound            = &CustomError{Code: ErrNotFoundCode, Message: "Not Found"}
	ErrConflict            = &CustomError{Code: ErrConflictCode, Message: "Conflict"}
	ErrUnauthorized        = &CustomError{Code: ErrUnauthorizedCode, Message: "Unauthorized"}
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
		default:
			return http.StatusInternalServerError
		}
	default:
		logrus.Warn("Received an unknown error type")
		return http.StatusInternalServerError
	}
}
