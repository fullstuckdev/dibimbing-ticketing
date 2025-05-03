package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/taufikmulyawan/ticketing-system/entity"
	"github.com/taufikmulyawan/ticketing-system/service"
)

// AuditMiddleware creates middleware for logging API access
func AuditMiddleware(auditService service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request first
		c.Next()
		
		// Skip logging for certain paths like health check, static files, etc.
		path := c.Request.URL.Path
		if path == "/health" || path == "/swagger/*any" {
			return
		}
		
		// Get user ID from context if available
		var userID uint
		if id, exists := c.Get("user_id"); exists {
			if parsedID, ok := id.(float64); ok {
				userID = uint(parsedID)
			}
		}
		
		// Skip audit logging if no user is authenticated
		if userID == 0 {
			return
		}
		
		// Get IP address and user agent
		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()
		
		// Determine action based on HTTP method
		var action string
		switch c.Request.Method {
		case "GET":
			action = "view"
		case "POST":
			action = "create"
		case "PUT", "PATCH":
			action = "update"
		case "DELETE":
			action = "delete"
		default:
			action = "other"
		}
		
		// Determine entity type based on URL path
		entityType := "unknown"
		if len(path) > 1 {
			segments := splitPath(path)
			if len(segments) > 0 {
				entityType = segments[0]
			}
		}
		
		// Log the activity (for simplicity, we're not capturing entity IDs, old values, or new values here)
		// In a real application, you might want to capture request/response bodies for more detailed logging
		go auditService.LogActivity(
			userID,
			entity.AuditAction(action),
			entityType,
			0, // entityID is not determined here
			nil, // oldValue
			nil, // newValue
			ipAddress,
			userAgent,
		)
	}
}

// Helper function to split path
func splitPath(path string) []string {
	var result []string
	var current string
	
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(path[i])
		}
	}
	
	if current != "" {
		result = append(result, current)
	}
	
	return result
} 