package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code     string      `json:"code"`
	Message  string      `json:"message"`
	Details  interface{} `json:"details,omitempty"`
	Metadata *ErrorMeta  `json:"metadata,omitempty"`
}

type ErrorMeta struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

type Meta struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func RespondSuccess(c *gin.Context, statusCode int, data interface{}, message string) {
	requestID, _ := c.Get("request_id")
	reqIDStr, _ := requestID.(string)

	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
		Meta: &Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			RequestID: reqIDStr,
		},
	})
}

func RespondError(c *gin.Context, statusCode int, code string, message string, details interface{}) {
	requestID, _ := c.Get("request_id")
	reqIDStr, _ := requestID.(string)

	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
			Metadata: &ErrorMeta{
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				RequestID: reqIDStr,
			},
		},
	})
}

func RespondPaginated(c *gin.Context, statusCode int, data interface{}, page, limit, total int) {
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	c.JSON(statusCode, PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
