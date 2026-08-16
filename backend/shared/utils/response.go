package utils

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type APIError struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Code    string            `json:"code,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
	Error   string            `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, status int, message string, data any) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, status int, message string) {
	ErrorResponseWithCode(c, status, "", message)
}

func ErrorResponseWithCode(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIError{
		Success: false,
		Message: message,
		Code:    code,
	})
}

func ValidationErrorResponse(c *gin.Context, status int, code, message string, errs map[string]string) {
	c.JSON(status, APIError{
		Success: false,
		Message: message,
		Code:    code,
		Errors:  errs,
	})
}
