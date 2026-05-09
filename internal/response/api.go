// Package response defines HTTP response builders.
// Layer 3: Depends on types (Layer 0). No business logic.
package response

import (
	"net/http"

	"github.com/example/dormitory-management/internal/types"
	"github.com/gin-gonic/gin"
)

// Response is the standard API response wrapper.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo contains error details.
type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Success returns a successful response.
func Success(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// Error returns an error response.
func Error(err error) Response {
	return Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_ERROR",
			Details: err.Error(),
		},
	}
}

// ErrorWithCode returns an error response with a specific code.
func ErrorWithCode(code int, message, details string) Response {
	return Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// Paginated returns a paginated response.
func Paginated(data interface{}, pagination types.Pagination) Response {
	return Success(gin.H{
		"data":       data,
		"pagination": pagination,
	})
}

// Created returns a 201 Created response.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Success(data))
}

// OK returns a 200 OK response.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Success(data))
}

// NoContent returns a 204 No Content response.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
