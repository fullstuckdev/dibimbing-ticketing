package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/taufikmulyawan/ticketing-system/service"
)

type AuditController interface {
	GetAuditLogs(c *gin.Context)
	GetEntityAuditLogs(c *gin.Context)
}

type auditController struct {
	auditService service.AuditService
}

func NewAuditController(auditService service.AuditService) AuditController {
	return &auditController{
		auditService: auditService,
	}
}

// GetAuditLogs godoc
// @Summary Get audit logs
// @Description Retrieve audit logs with optional filtering by user, entity type, and date range
// @Tags audit
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Param user_id query int false "Filter by user ID"
// @Param entity_type query string false "Filter by entity type (e.g., 'user', 'event', 'ticket')"
// @Param start_date query string false "Start date filter (format: YYYY-MM-DD)"
// @Param end_date query string false "End date filter (format: YYYY-MM-DD)"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of audit logs with pagination info"
// @Router /audit/logs [get]
func (ctrl *auditController) GetAuditLogs(c *gin.Context) {
	// Parse query parameters with defaults
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	userIDStr := c.Query("user_id")
	entityType := c.Query("entity_type")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	
	// Convert userID to uint
	var userID uint
	if userIDStr != "" {
		if id, err := strconv.Atoi(userIDStr); err == nil {
			userID = uint(id)
		}
	}
	
	// Parse date strings to time.Time
	var startDate, endDate time.Time
	var err error
	
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
			return
		}
	}
	
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
			return
		}
		// Set end date to the end of the day
		endDate = endDate.Add(24*time.Hour - time.Second)
	}
	
	// Get audit logs based on filters
	logs, total, err := ctrl.auditService.GetAuditLogs(page, limit, userID, entityType, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audit logs"})
		return
	}
	
	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)
	hasNext := int64(page) < totalPages
	hasPrev := page > 1
	
	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    hasNext,
			"has_prev":    hasPrev,
		},
	})
}

// GetEntityAuditLogs godoc
// @Summary Get audit logs for a specific entity
// @Description Retrieve audit logs for a specific entity by type and ID
// @Tags audit
// @Accept json
// @Produce json
// @Param entity_type path string true "Entity type (e.g., 'user', 'event', 'ticket')"
// @Param entity_id path int true "Entity ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of audit logs for the entity"
// @Router /audit/{entity_type}/{entity_id} [get]
func (ctrl *auditController) GetEntityAuditLogs(c *gin.Context) {
	entityType := c.Param("entity_type")
	entityIDStr := c.Param("entity_id")
	
	// Convert entityID to uint
	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}
	
	// Get audit logs for the entity
	logs, err := ctrl.auditService.GetAuditLogsByEntity(entityType, uint(entityID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audit logs"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": logs})
} 