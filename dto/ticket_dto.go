package dto

import (
	"time"
)

// TicketCreateRequest represents the request for creating a new ticket
type TicketCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	EventID     uint   `json:"event_id" binding:"required"`
}

// TicketUpdateRequest represents the request for updating a ticket
type TicketUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// TicketResponse represents the response format for ticket data
type TicketResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	UserID      uint      `json:"user_id"`
	EventID     uint      `json:"event_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TicketDetailResponse represents detailed ticket data with related information
type TicketDetailResponse struct {
	TicketResponse
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	EventName string `json:"event_name"`
}

// TicketFilterRequest represents filters for ticket listing
type TicketFilterRequest struct {
	Status    string `form:"status"`
	UserID    uint   `form:"user_id"`
	EventID   uint   `form:"event_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page,default=1"`
	Limit     int    `form:"limit,default=10"`
}

// TicketListResponse represents paginated list of tickets
type TicketListResponse struct {
	Data       []TicketResponse `json:"data"`
	Pagination PaginationMeta   `json:"meta"`
} 