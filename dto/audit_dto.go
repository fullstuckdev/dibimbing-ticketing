package dto

import (
	"time"
)

// AuditLogResponse represents the response format for audit log data
type AuditLogResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	UserName   string    `json:"user_name"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   uint      `json:"entity_id"`
	OldData    string    `json:"old_data,omitempty"`
	NewData    string    `json:"new_data,omitempty"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}

// AuditLogFilterRequest represents filters for audit log listing
type AuditLogFilterRequest struct {
	UserID     uint   `form:"user_id"`
	Action     string `form:"action"`
	EntityType string `form:"entity_type"`
	EntityID   uint   `form:"entity_id"`
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=10"`
}

// AuditLogListResponse represents paginated list of audit logs
type AuditLogListResponse struct {
	Data       []AuditLogResponse `json:"data"`
	Pagination PaginationMeta     `json:"meta"`
} 