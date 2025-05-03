package service

import (
	"errors"
	"time"

	"github.com/taufikmulyawan/ticketing-system/entity"
	"github.com/taufikmulyawan/ticketing-system/repository"
)

type EventService interface {
	GetAllEvents(page, limit int) ([]entity.Event, int64, error)
	GetEventByID(id uint) (*entity.Event, error)
	CreateEvent(event *entity.Event) error
	UpdateEvent(id uint, event *entity.Event) error
	DeleteEvent(id uint) error
}

type eventService struct {
	eventRepo repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{
		eventRepo: eventRepo,
	}
}

func (s *eventService) GetAllEvents(page, limit int) ([]entity.Event, int64, error) {
	// Default pagination values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	return s.eventRepo.FindAll(page, limit)
}

func (s *eventService) GetEventByID(id uint) (*entity.Event, error) {
	return s.eventRepo.FindByID(id)
}

func (s *eventService) CreateEvent(event *entity.Event) error {
	// Validate event data
	if event.Name == "" {
		return errors.New("event name is required")
	}
	if event.Location == "" {
		return errors.New("event location is required")
	}
	if event.Capacity <= 0 {
		return errors.New("event capacity must be positive")
	}
	if event.Price < 0 {
		return errors.New("event price cannot be negative")
	}
	if event.StartDate.Before(time.Now()) {
		return errors.New("event start date must be in the future")
	}
	if event.EndDate.Before(event.StartDate) {
		return errors.New("event end date must be after start date")
	}

	// Set default status
	if event.Status == "" {
		event.Status = entity.EventStatusActive
	}

	// Save event
	return s.eventRepo.Save(event)
}

func (s *eventService) UpdateEvent(id uint, event *entity.Event) error {
	// Get existing event
	existingEvent, err := s.eventRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check if event is already finished
	if existingEvent.Status == entity.EventStatusFinished {
		return errors.New("cannot update a finished event")
	}

	// Update event fields
	existingEvent.Name = event.Name
	existingEvent.Description = event.Description
	existingEvent.Location = event.Location
	existingEvent.StartDate = event.StartDate
	existingEvent.EndDate = event.EndDate
	existingEvent.Capacity = event.Capacity
	existingEvent.Price = event.Price
	existingEvent.Status = event.Status

	// Save updated event
	return s.eventRepo.Save(existingEvent)
}

func (s *eventService) DeleteEvent(id uint) error {
	return s.eventRepo.Delete(id)
} 