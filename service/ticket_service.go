package service

import (
	"errors"
	"time"

	"github.com/taufikmulyawan/ticketing-system/entity"
	"github.com/taufikmulyawan/ticketing-system/repository"
)

type TicketService interface {
	GetAllTickets(page, limit int, userID uint) ([]entity.Ticket, int64, error)
	GetTicketByID(id uint) (*entity.Ticket, error)
	PurchaseTicket(ticket *entity.Ticket) error
	CancelTicket(id uint, userID uint) error
}

type ticketService struct {
	ticketRepo repository.TicketRepository
	eventRepo  repository.EventRepository
}

func NewTicketService(ticketRepo repository.TicketRepository, eventRepo repository.EventRepository) TicketService {
	return &ticketService{
		ticketRepo: ticketRepo,
		eventRepo:  eventRepo,
	}
}

func (s *ticketService) GetAllTickets(page, limit int, userID uint) ([]entity.Ticket, int64, error) {
	// Default pagination values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	return s.ticketRepo.FindAll(page, limit, userID)
}

func (s *ticketService) GetTicketByID(id uint) (*entity.Ticket, error) {
	return s.ticketRepo.FindByID(id)
}

func (s *ticketService) PurchaseTicket(ticket *entity.Ticket) error {
	// Check if event exists
	event, err := s.eventRepo.FindByID(ticket.EventID)
	if err != nil {
		return err
	}

	// Check if event is active
	if event.Status != entity.EventStatusActive {
		return errors.New("tickets can only be purchased for active events")
	}

	// Check if event date is in the future
	if event.StartDate.Before(time.Now()) {
		return errors.New("cannot purchase tickets for past events")
	}

	// Check if event has available capacity
	soldTickets, err := s.ticketRepo.CountSoldTicketsByEventID(event.ID)
	if err != nil {
		return err
	}

	if int(soldTickets) >= event.Capacity {
		return errors.New("event is sold out")
	}

	// Set ticket details
	ticket.Status = entity.TicketStatusPurchased
	ticket.PurchasedAt = time.Now()

	// Save the ticket
	return s.ticketRepo.Save(ticket)
}

func (s *ticketService) CancelTicket(id uint, userID uint) error {
	// Get the ticket
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check if the ticket belongs to the user (unless admin)
	if ticket.UserID != userID {
		return errors.New("unauthorized to cancel this ticket")
	}

	// Check if the ticket is already cancelled
	if ticket.Status == entity.TicketStatusCancelled {
		return errors.New("ticket is already cancelled")
	}

	// Check if the event has already started
	if ticket.Event.StartDate.Before(time.Now()) {
		return errors.New("cannot cancel tickets for events that have already started")
	}

	// Update the ticket status
	ticket.Status = entity.TicketStatusCancelled
	return s.ticketRepo.Save(ticket)
} 