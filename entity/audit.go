package entity

import (
	"time"
)

// AuditAction represents the type of action performed
type AuditAction string

const (
	ActionCreate AuditAction = "create"
	ActionUpdate AuditAction = "update"
	ActionDelete AuditAction = "delete"
	ActionLogin  AuditAction = "login"
	ActionLogout AuditAction = "logout"
)

// AuditLog represents an audit trail entry in the system
type AuditLog struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `json:"user_id"`
	Action    AuditAction `gorm:"size:50;not null" json:"action"`
	EntityType string     `gorm:"size:50;not null" json:"entity_type"` // e.g., "user", "event", "ticket"
	EntityID   uint       `json:"entity_id"`
	OldValue   string     `gorm:"type:text" json:"old_value,omitempty"`
	NewValue   string     `gorm:"type:text" json:"new_value,omitempty"`
	IPAddress  string     `gorm:"size:50" json:"ip_address,omitempty"`
	UserAgent  string     `gorm:"size:255" json:"user_agent,omitempty"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	
	// Navigation property
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
} 