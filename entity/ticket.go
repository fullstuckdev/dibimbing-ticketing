package entity

import (
	"time"
)

type TicketStatus string

const (
	TicketStatusAvailable  TicketStatus = "available"
	TicketStatusPurchased  TicketStatus = "purchased"
	TicketStatusCancelled  TicketStatus = "cancelled"
)

type Ticket struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	UserID    uint         `gorm:"not null" json:"user_id"`
	EventID   uint         `gorm:"not null" json:"event_id"`
	Status    TicketStatus `gorm:"size:50;not null;default:purchased" json:"status"`
	PurchasedAt time.Time  `gorm:"not null" json:"purchased_at"`
	CreatedAt time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	User      User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Event     Event        `gorm:"foreignKey:EventID" json:"event,omitempty"`
} 