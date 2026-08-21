package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func ResponseSuccess(c *gin.Context, httpCode int, message string, data interface{}) {
	c.JSON(httpCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ResponseError(c *gin.Context, httpCode int, message string, errors interface{}) {
	c.JSON(httpCode, APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func InternalServerError(c *gin.Context, message string) {
	ResponseError(c, http.StatusInternalServerError, message, nil)
}
