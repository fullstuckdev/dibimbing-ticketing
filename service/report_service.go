package service

import (
	"github.com/taufikmulyawan/ticketing-system/reports"
	"github.com/taufikmulyawan/ticketing-system/repository"
	"github.com/taufikmulyawan/ticketing-system/types"
)

// Type aliases for backward compatibility
type EventSalesSummary = types.EventSalesSummary
type SalesSummary = types.SalesSummary

type ReportService interface {
	GetSalesSummary() (*SalesSummary, error)
	GetEventSalesSummary(eventID uint) (*EventSalesSummary, error)
	ExportSalesSummaryPDF() ([]byte, error) 
	ExportEventSalesPDF(eventID uint) ([]byte, error)
	ExportSalesSummaryCSV() ([]byte, error)
	ExportEventSalesCSV(eventID uint) ([]byte, error)
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

// Export methods

func (s *reportService) ExportSalesSummaryPDF() ([]byte, error) {
	summary, err := s.GetSalesSummary()
	if err != nil {
		return nil, err
	}
	
	return reports.GenerateSalesSummaryPDF(summary)
}

func (s *reportService) ExportEventSalesPDF(eventID uint) ([]byte, error) {
	summary, err := s.GetEventSalesSummary(eventID)
	if err != nil {
		return nil, err
	}
	
	return reports.GenerateEventSalesPDF(summary)
}

func (s *reportService) ExportSalesSummaryCSV() ([]byte, error) {
	summary, err := s.GetSalesSummary()
	if err != nil {
		return nil, err
	}
	
	return reports.GenerateSalesSummaryCSV(summary)
}

func (s *reportService) ExportEventSalesCSV(eventID uint) ([]byte, error) {
	summary, err := s.GetEventSalesSummary(eventID)
	if err != nil {
		return nil, err
	}
	
	return reports.GenerateEventSalesCSV(summary)
} 