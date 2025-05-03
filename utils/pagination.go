package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPaginationParams extracts pagination parameters from the request
func GetPaginationParams(c *gin.Context) (page int, limit int) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Cap the limit to prevent large queries
	if limit > 100 {
		limit = 100
	}

	return page, limit
}

// GeneratePaginationResponse creates a standardized pagination response
func GeneratePaginationResponse(data interface{}, page int, limit int, total int64) gin.H {
	totalPages := (total + int64(limit) - 1) / int64(limit)

	return gin.H{
		"data": data,
		"meta": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total_items":  total,
			"total_pages":  totalPages,
		},
	}
} 