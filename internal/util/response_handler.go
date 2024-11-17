package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status     bool        `json:"status"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:     true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func SendError(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:     false,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}
