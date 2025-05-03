package service

import (
	"github.com/taufikmulyawan/ticketing-system/repository"
)

type EventSalesSummary struct {
	EventID       uint    `json:"event_id"`
	EventName     string  `json:"event_name"`
	TotalTickets  int64   `json:"total_tickets"`
	TotalRevenue  float64 `json:"total_revenue"`
}

type SalesSummary struct {
	TotalEvents   int64              `json:"total_events"`
	TotalTickets  int64              `json:"total_tickets"`
	TotalRevenue  float64            `json:"total_revenue"`
	EventSummary  []EventSalesSummary `json:"event_summary"`
}

type ReportService interface {
	GetSalesSummary() (*SalesSummary, error)
	GetEventSalesSummary(eventID uint) (*EventSalesSummary, error)
}

type reportService struct {
	ticketRepo repository.TicketRepository
	eventRepo  repository.EventRepository
}

func NewReportService(ticketRepo repository.TicketRepository, eventRepo repository.EventRepository) ReportService {
	return &reportService{
		ticketRepo: ticketRepo,
		eventRepo:  eventRepo,
	}
}

func (s *reportService) GetSalesSummary() (*SalesSummary, error) {
	// Get all events for calculating summary
	events, _, err := s.eventRepo.FindAll(1, 1000) // Using large limit to get all events
	if err != nil {
		return nil, err
	}

	summary := &SalesSummary{
		TotalEvents:  int64(len(events)),
		TotalTickets: 0,
		TotalRevenue: 0,
		EventSummary: make([]EventSalesSummary, 0),
	}

	// Calculate summary for each event
	for _, event := range events {
		soldTickets, err := s.ticketRepo.CountSoldTicketsByEventID(event.ID)
		if err != nil {
			return nil, err
		}

		revenue := float64(soldTickets) * event.Price
		summary.TotalTickets += soldTickets
		summary.TotalRevenue += revenue

		eventSummary := EventSalesSummary{
			EventID:      event.ID,
			EventName:    event.Name,
			TotalTickets: soldTickets,
			TotalRevenue: revenue,
		}

		summary.EventSummary = append(summary.EventSummary, eventSummary)
	}

	return summary, nil
}

func (s *reportService) GetEventSalesSummary(eventID uint) (*EventSalesSummary, error) {
	// Get the event
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, err
	}

	// Count sold tickets for this event
	soldTickets, err := s.ticketRepo.CountSoldTicketsByEventID(event.ID)
	if err != nil {
		return nil, err
	}

	revenue := float64(soldTickets) * event.Price
	eventSummary := &EventSalesSummary{
		EventID:      event.ID,
		EventName:    event.Name,
		TotalTickets: soldTickets,
		TotalRevenue: revenue,
	}

	return eventSummary, nil
} 