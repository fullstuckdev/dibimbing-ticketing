package types

// EventSalesSummary represents sales data for a specific event
type EventSalesSummary struct {
	EventID       uint    `json:"event_id"`
	EventName     string  `json:"event_name"`
	TotalTickets  int64   `json:"total_tickets"`
	TotalRevenue  float64 `json:"total_revenue"`
}

// SalesSummary represents overall sales data across all events
type SalesSummary struct {
	TotalEvents   int64              `json:"total_events"`
	TotalTickets  int64              `json:"total_tickets"`
	TotalRevenue  float64            `json:"total_revenue"`
	EventSummary  []EventSalesSummary `json:"event_summary"`
} 