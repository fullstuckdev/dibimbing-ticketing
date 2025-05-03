package entity

import (
	"time"
)

type EventStatus string

const (
	EventStatusActive   EventStatus = "active"
	EventStatusOngoing  EventStatus = "ongoing"
	EventStatusFinished EventStatus = "finished"
)

type Event struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Name        string      `gorm:"size:255;not null;unique" json:"name"`
	Description string      `gorm:"type:text" json:"description"`
	Location    string      `gorm:"size:255;not null" json:"location"`
	StartDate   time.Time   `gorm:"not null" json:"start_date"`
	EndDate     time.Time   `gorm:"not null" json:"end_date"`
	Capacity    int         `gorm:"not null" json:"capacity"`
	Price       float64     `gorm:"not null" json:"price"`
	Status      EventStatus `gorm:"size:50;not null;default:active" json:"status"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	Tickets     []Ticket    `gorm:"foreignKey:EventID" json:"tickets,omitempty"`
} 