package util

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	ErrInternalServerError = "INTERNAL_SERVER_ERROR"
	ErrNotFound            = "NOT_FOUND"
	ErrConflict            = "CONFLICT"
)

type CustomError struct {
	Code    string
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)

	switch e := err.(type) {
	case *CustomError:
		switch e.Code {
		case ErrInternalServerError:
			return http.StatusInternalServerError
		case ErrNotFound:
			return http.StatusNotFound
		case ErrConflict:
			return http.StatusConflict
		default:
			return http.StatusInternalServerError
		}
	default:
		return http.StatusInternalServerError
	}
}
