package dto

import (
	"time"
)

// EventCreateRequest represents the request for creating a new event
type EventCreateRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Location    string    `json:"location"`
	Capacity    int       `json:"capacity"`
}

// EventUpdateRequest represents the request for updating an event
type EventUpdateRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Location    string    `json:"location"`
	Capacity    int       `json:"capacity"`
}

// EventResponse represents the response format for event data
type EventResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Location    string    `json:"location"`
	Capacity    int       `json:"capacity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EventDetailResponse represents detailed event data with ticket information
type EventDetailResponse struct {
	EventResponse
	TicketCount    int  `json:"ticket_count"`
	AvailableSeats int  `json:"available_seats"`
	IsActive       bool `json:"is_active"`
}

// EventFilterRequest represents filters for event listing
type EventFilterRequest struct {
	Name      string `form:"name"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Location  string `form:"location"`
	Page      int    `form:"page,default=1"`
	Limit     int    `form:"limit,default=10"`
}

// EventListResponse represents paginated list of events
type EventListResponse struct {
	Data       []EventResponse `json:"data"`
	Pagination PaginationMeta  `json:"meta"`
} 