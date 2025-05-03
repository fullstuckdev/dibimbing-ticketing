package middleware

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/taufikmulyawan/ticketing-system/entity"
	"github.com/taufikmulyawan/ticketing-system/service"
)

// AuditMiddleware creates middleware for logging API access
func AuditMiddleware(auditService service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Store the request body for audit purposes
		var requestBody []byte
		var err error
		
		// For POST, PUT, PATCH methods, we want to capture the request body
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			// Read the request body
			requestBody, err = ioutil.ReadAll(c.Request.Body)
			if err == nil {
				// Restore the request body so it can be read again by handlers
				c.Request.Body = ioutil.NopCloser(strings.NewReader(string(requestBody)))
			}
		}
		
		// Process request first
		c.Next()
		
		// Skip logging for certain paths like health check, static files, etc.
		path := c.Request.URL.Path
		if path == "/health" || strings.HasPrefix(path, "/swagger/") || path == "/favicon.ico" {
			return
		}
		
		// Get user ID from context if available
		var userID uint
		if id, exists := c.Get("user_id"); exists {
			if parsedID, ok := id.(float64); ok {
				userID = uint(parsedID)
			} else if parsedID, ok := id.(uint); ok {
				userID = parsedID
			}
		}
		
		// Set userID to 0 for public routes like login/register
		// We'll still log these actions but without a user association
		
		// Get IP address and user agent
		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()
		
		// Determine action based on HTTP method
		var action entity.AuditAction
		switch c.Request.Method {
		case "GET":
			action = "view"
		case "POST":
			if path == "/login" {
				action = entity.ActionLogin
			} else if path == "/register" {
				action = entity.ActionCreate
			} else {
				action = entity.ActionCreate
			}
		case "PUT", "PATCH":
			action = entity.ActionUpdate
		case "DELETE":
			action = entity.ActionDelete
		default:
			action = "other"
		}
		
		// Determine entity type and ID based on URL path
		entityType, entityID := extractEntityInfo(c)
		
		// Prepare oldValue and newValue based on request and response
		var oldValue, newValue interface{}
		
		// For POST, PUT, PATCH, capture the request body as newValue
		if (c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH") && len(requestBody) > 0 {
			// Parse the request body (could be JSON, form data, etc.)
			var bodyData interface{}
			if err := json.Unmarshal(requestBody, &bodyData); err == nil {
				newValue = bodyData
			} else {
				// If not valid JSON, just store as string
				newValue = string(requestBody)
			}
		}
		
		// Log the activity
		go auditService.LogActivity(
			userID,
			action,
			entityType,
			entityID,
			oldValue,
			newValue,
			ipAddress,
			userAgent,
		)
	}
}

// Helper function to extract entity type and ID from the URL
func extractEntityInfo(c *gin.Context) (string, uint) {
	path := c.Request.URL.Path
	segments := splitPath(path)
	
	// Default values
	entityType := "unknown"
	var entityID uint
	
	if len(segments) > 0 {
		// First segment is usually the entity type (users, events, tickets, etc.)
		entityType = segments[0]
		
		// If we have an ID in the URL (like /events/123), extract it
		if len(segments) > 1 {
			if id, err := strconv.Atoi(segments[1]); err == nil {
				entityID = uint(id)
			}
		}
	}
	
	// Special case for login and register
	if path == "/login" {
		entityType = "auth"
	} else if path == "/register" {
		entityType = "user"
	}
	
	return entityType, entityID
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